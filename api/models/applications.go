package models

import (
	"errors"
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/config"
	"github.com/lucasmbaia/forcloudy/api/core-xmpp"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/repository"
	"github.com/lucasmbaia/forcloudy/api/utils"
	"github.com/satori/go.uuid"
	"strconv"
)

const (
	PATH_DEFAULT  = "/root/go/src/github.com/lucasmbaia/forcloudy/minion/core/"
	BUILD_DEFAULT = "hello_world"
)

type Applications struct {
	repository repository.Repositorier
}

func NewApplications(session repository.Repositorier) *Applications {
	return &Applications{repository: session}
}

type ApplicationEtcd struct {
	Protocol        map[string]string `json:"protocol,omitempty"`
	Image           string            `json:"image,omitempty"`
	PortsDST        []string          `json:"portsDst,omitempty"`
	Cpus            string            `json:"cpus,omitempty"`
	Dns             string            `json:"dns,omitempty"`
	Memory          string            `json:"memory,omitempty"`
	TotalContainers int               `json:"totalContainers,omitempty"`
}

func (a *Applications) Post(values interface{}) error {
	var (
		application  = values.(*datamodels.ApplicationsFields)
		applications interface{}
		err          error
		customers    interface{}
		image        string
		imageID      uuid.UUID
		containerID  uuid.UUID
		iterator     = 1
	)

	if applications, err = a.Get(datamodels.ApplicationsFields{Name: application.Name, Customer: application.Customer}); err != nil {
		return err
	}

	if len(applications.([]datamodels.ApplicationsFields)) > 0 {
		return errors.New(fmt.Sprintf("Name of application %s exists in database", application.Name))
	}

	if customers, err = NewCustomers(a.repository).Get(
		datamodels.CustomersFields{
			ID: application.Customer,
		},
	); err != nil {
		return err
	} else {
		if len(customers.([]datamodels.CustomersFields)) == 0 {
			return errors.New("Invalid Customer")
		}
	}

	image = fmt.Sprintf("%s_app-%s", customers.([]datamodels.CustomersFields)[0].Name, application.Name)
	if imageID, err = uuid.NewV4(); err != nil {
		return err
	}

	if err = NewImages(a.repository).Post(
		&datamodels.ImagesFields{
			ID:       imageID.String(),
			Customer: application.Customer,
			Name:     image,
			Version:  "v1",
		},
	); err != nil {
		return err
	}

	application.Image = imageID.String()
	application.Status = "IN_PROGRESS"

	if err = a.repository.Create(application); err != nil {
		return err
	}

	for iterator <= application.TotalContainers {
		if containerID, err = uuid.NewV4(); err != nil {
			return err
		}

		if err = NewContainers(a.repository).Post(
			&datamodels.ContainersFields{
				ID:          containerID.String(),
				Customer:    application.Customer,
				Application: application.ID,
				Name:        fmt.Sprintf("%s_app-%s-%d", customers.([]datamodels.CustomersFields)[0].Name, application.Name, iterator),
				Status:      "IN_PROGRESS",
				State:       "CREATING",
			},
		); err != nil {
			return err
		}

		iterator++
	}

	go a.requestDeploy(application, customers.([]datamodels.CustomersFields)[0].Name)
	return nil
}

func (a *Applications) Get(filters interface{}) (interface{}, error) {
	var (
		entity = []datamodels.ApplicationsFields{}
		err    error
	)

	if _, err = a.repository.Read(filters, &entity); err != nil {
		return entity, err
	}

	return entity, err
}

func (a *Applications) Delete(conditions interface{}) error {
	return nil
}

func (a *Applications) Put(fields, data interface{}) error {
	return nil
}

func (a *Applications) Patch(fields, data interface{}) error {
	var (
		conditions = fields.(*datamodels.ApplicationsFields)
		entity     = data.(*datamodels.ApplicationsFields)
	)

	return a.repository.Update(conditions, entity)
}

func (a *Applications) requestDeploy(application *datamodels.ApplicationsFields, customer string) {
	var (
		response        = make(chan core.Container)
		ports           []core.Ports
		err             error
		entity          = &datamodels.ApplicationsFields{Status: "COMPLETED"}
		applicationEtcd ApplicationEtcd
		key             string
	)

	go func() {
		for {
			select {
			case resp := <-response:
				fmt.Println(resp)
				var (
					data           = &datamodels.ContainersFields{Status: "COMPLETED", State: "CREATED"}
					portsContainer = make(map[string][]string)
					protocol       = make(map[string]string)
				)

				if resp.Error != nil {
					data.Status = "ERROR"
					data.State = "ERROR"
				}

				if err = NewContainers(a.repository).Patch(
					&datamodels.ContainersFields{
						Application: application.ID,
						Name:        resp.Name,
					},
					data,
				); err != nil {
					fmt.Println(err)
				}

				for _, port := range resp.PortsContainer {
					portsContainer[port.Source] = port.Destinations
				}

				for _, port := range application.Ports {
					protocol[strconv.Itoa(port.Port)] = port.Protocol
				}

				if err = utils.GenerateConf(utils.Haproxy{
					Customer:         customer,
					ApplicationName:  application.Name,
					ContainerName:    resp.Name,
					PortsContainer:   portsContainer,
					Protocol:         protocol,
					AddressContainer: resp.Address,
					Dns:              application.Dns,
					Minion:           resp.Minion,
				}); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	applicationEtcd = ApplicationEtcd{
		Image:           fmt.Sprintf("%s_app-%s/image:v1", customer, application.Name),
		Cpus:            application.Cpus,
		Memory:          fmt.Sprintf("%dMB", ((application.Memory / 1024) / 1024)),
		TotalContainers: application.TotalContainers,
		Protocol:        make(map[string]string),
	}

	for _, port := range application.Ports {
		ports = append(ports, core.Ports{Port: port.Port, Protocol: port.Protocol})
		applicationEtcd.Protocol[strconv.Itoa(port.Port)] = port.Protocol
		applicationEtcd.PortsDST = append(applicationEtcd.PortsDST, strconv.Itoa(port.Port))
	}

	key = fmt.Sprintf("/%s/%s", customer, application.Name)

	if err = config.EnvSingleton.EtcdConnection.Set(key, applicationEtcd); err != nil {
		fmt.Println(err)
	}

	if err = core.DeployApplication(core.Deploy{
		Customer:        customer,
		ApplicationName: application.Name,
		ImageVersion:    "v1",
		Cpus:            application.Cpus,
		Memory:          fmt.Sprintf("%dMB", ((application.Memory / 1024) / 1024)),
		TotalContainers: application.TotalContainers,
		Ports:           ports,
		Path:            PATH_DEFAULT,
		Build:           BUILD_DEFAULT,
	}, 1, true, response); err != nil {
		*entity.Error = err.Error()
	}

	if err = a.Patch(
		&datamodels.ApplicationsFields{
			ID: application.ID,
		},
		entity,
	); err != nil {
		fmt.Println(err)
	}

	return
}
