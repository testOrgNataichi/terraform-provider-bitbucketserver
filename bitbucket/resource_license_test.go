package bitbucket

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketLicense(t *testing.T) {
	testAccBitbucketLicenseConfig := `
		resource "bitbucketserver_license" "test" {
			license = "AAABrQ0ODAoPeNp9kVFvmzAQx9/9KSztLZIJZIu0RkJqA6yNViAK0G3d+uDApXgjNrKPbPn2dYG0a 6fuwS8+393v//O7vAMa8yN1PerOFrPZYn5GL+OczlzvI0m6/RZ0uisMaON7LgmURF5iwvfgVy3XW pj6nGPDjRFcOqXaE4Pc1M61KEEayI8t9I+DNI6jTbC6uP73wd/FdafLmhsIOYL/yMDcOXM98p95Y yn60wp97PvW769OpFHMRfMWagb6AHoV+svLs5x9LW4+sM+3t1ds6XpfRkw7jwcgEbSPugOSdVtTa tGiUHK4mUwmSZqzT+mGrTdpWAT5Kk1YkUW24AcaLFBFt0eKNdARlUayVBVo2mr1E0qk32vE9sdiO r1XzgvEaTN0MBg67hwaKioV0koY1GLbIdjJwlBUtOwMqr39KYfY1JZZclm+9jLEsmbEAZ4CBJvoI o9Ctvz2CP2GrRHe6irkL6l+S5JFiW8Pm7suSfU9l8LwXkwIB2hUaxPmYPAUm/Q2bP315w5MGXL95 DmEZ839jFEE3SlNedvS6rTCkOjAm25YvOON3fMAVTj4nTAtAhRH4o+fI5MQ7xSh2mtA1bPJrq0WA gIVAIGperR8m2N0fl/GfUUJfQnd+T1aX02kk"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketLicenseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketLicenseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketLicenseExists("bitbucketserver_license.test"),
				),
			},
		},
	})
}

func testAccCheckBitbucketLicenseDestroy(s *terraform.State) error {
	_, ok := s.RootModule().Resources["bitbucketserver_license.test"]
	if !ok {
		return fmt.Errorf("not found %s", "bitbucketserver_license.test")
	}

	return nil
}

func testAccCheckBitbucketLicenseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no license ID is set")
		}
		return nil
	}
}
