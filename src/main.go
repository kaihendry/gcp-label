package label

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type pubSubMessage struct {
	Data []byte `json:"data"`
}

type instanceEventTemplate struct {
	InsertID  string `json:"insertId"`
	LogName   string `json:"logName"`
	Operation struct {
		ID       string `json:"id"`
		Last     bool   `json:"last"`
		Producer string `json:"producer"`
	} `json:"operation"`
	ProtoPayload struct {
		Type               string `json:"@type"`
		AuthenticationInfo struct {
			PrincipalEmail string `json:"principalEmail"`
		} `json:"authenticationInfo"`
		MethodName string `json:"methodName"`
		Request    struct {
			Type string `json:"@type"`
		} `json:"request"`
		RequestMetadata struct {
			CallerSuppliedUserAgent string `json:"callerSuppliedUserAgent"`
		} `json:"requestMetadata"`
		ResourceName string `json:"resourceName"`
		ServiceName  string `json:"serviceName"`
	} `json:"protoPayload"`
	ReceiveTimestamp time.Time `json:"receiveTimestamp"`
	Resource         struct {
		Labels struct {
			InstanceID string `json:"instance_id"`
			ProjectID  string `json:"project_id"`
			Zone       string `json:"zone"`
		} `json:"labels"`
		Type string `json:"type"`
	} `json:"resource"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
}

// Label consumes a Pub/Sub message.
func Label(ctx context.Context, m pubSubMessage) (err error) {

	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return
	}
	log.Println("Project ID", credentials.ProjectID)

	var instanceEvent instanceEventTemplate
	err = json.Unmarshal(m.Data, &instanceEvent)
	if err != nil {
		return
	}
	log.Println(instanceEvent)

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return
	}

	computeService, err := compute.New(c)
	if err != nil {
		return
	}

	s := strings.Split(instanceEvent.ProtoPayload.ResourceName, "/")
	project, zone, instance := s[1], s[3], s[5]
	log.Printf("Split %s, to project: %s, zone: %s, instance: %s",
		s,
		project,
		zone,
		instance)

	inst, err := computeService.Instances.Get(project, zone, instance).Do()
	if err != nil {
		return
	}

	existingLabels := inst.Labels
	if existingLabels == nil {
		existingLabels = make(map[string]string)
	}

	existingLabels["instance-name"] = instance
	existingLabels["instance-id"] = instanceEvent.Resource.Labels.InstanceID

	log.Println("Applying labels:", existingLabels)

	rb := &compute.InstancesSetLabelsRequest{
		LabelFingerprint: inst.LabelFingerprint,
		Labels:           existingLabels,
	}

	// https://cloud.google.com/compute/docs/reference/rest/v1/instances/setLabels
	_, err = computeService.Instances.SetLabels(project, zone, instance, rb).Context(ctx).Do()
	if err != nil {
		return
	}

	return
}
