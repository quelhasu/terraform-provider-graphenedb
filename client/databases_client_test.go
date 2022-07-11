package client_test

import (
	"reflect"

	. "github.com/quelhasu/terraform-provider-graphenedb/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockResourceClient struct {
}

func (c *MockResourceClient) CreateResource(requestBody interface{}, responseBody interface{}) error {
	return nil
}

var _ = Describe("DatabasesClient", func() {
	var client DatabasesClient
	var authClient *AuthenticatedClient
	var resourceClient ResourceClient

	BeforeEach(func() {
		authClient = &AuthenticatedClient{}
		resourceClient = &MockResourceClient{}
		client = *authClient.NewDatabasesClient(resourceClient)
	})

	Describe("Initialization of client", func() {
		Context("a resource client is provided", func() {
			It("should be composed with injected ResourceClient", func() {
				Expect(client.ResourceClient).To(Equal(resourceClient))
			})
		})

		Context("a resource client is omitted", func() {
			BeforeEach(func() {
				client = *authClient.NewDatabasesClient()
			})

			It("should be composed with the DefaultResourceClient", func() {
				Expect(client.ResourceClient).ToNot(BeNil())
				Expect(reflect.TypeOf(client.ResourceClient).String()).To(Equal("*client.DefaultResourceClient"))
			})

			It("should be composed with correct ResourceRootPath", func() {
				Expect(client.ResourceClient.(*DefaultResourceClient).ResourceRootPath).To(Equal("/databases"))
			})

			It("should be composed with correct ResourceDescription", func() {
				Expect(client.ResourceClient.(*DefaultResourceClient).ResourceDescription).To(Equal("database list"))
			})

			It("should be composed with correct AuthenticatedClient", func() {
				Expect(client.ResourceClient.(*DefaultResourceClient).AuthenticatedClient).To(Equal(authClient))
			})
		})
	})

	Describe("Create database", func() {
		name := "db_name"
		version := "db_version"
		awsRegion := "db_awsRegion"
		plan := "db_plan"
		cidr := "db_cidr"

		Context("returns DatabaseDetail", func() {
			It("should set resource root path", func() {
				client.CreateDatabase(name, version, awsRegion, plan, cidr)
				Expect(client).ToNot(BeNil())
			})
		})
		// ID        string `json:"id"`
		// Name      string `json:"name"`
		// Version   string `json:"version"`
		// AwsRegion string `json:"awsRegion"`
		// Plan      string `json:"plan"`
		// URI       string `json:"uri"`
	})
})
