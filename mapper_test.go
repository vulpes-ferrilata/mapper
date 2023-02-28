package mapper_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vulpes-ferrilata/mapper"
)

var _ = Describe("Mapper", func() {
	type User struct {
		Name string
	}

	type Staff struct {
		Name string
	}

	When("mapping function was not registered", func() {
		It("cannot map", func() {
			user := &User{
				Name: "Vulpes Ferrilata",
			}

			staff, err := mapper.Map[User, Staff](user)

			Expect(staff).Should(BeNil())
			Expect(err).Should(MatchError(mapper.ErrMappingFunctionWasNotRegistered))
		})
	})

	When("mapping function was registered", func() {
		BeforeEach(func() {
			mapper.CreateMap(func(user *User) (*Staff, error) {
				if user == nil {
					return nil, nil
				}

				return &Staff{
					Name: user.Name,
				}, nil
			})
		})

		It("can map", func() {
			user := &User{
				Name: "Vulpes Ferrilata",
			}

			staff, err := mapper.Map[User, Staff](user)

			Expect(staff.Name).Should(BeEquivalentTo(user.Name))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("mapping function with custom error was registered", func() {
		userNotFoundErr := errors.New("user not found")

		BeforeEach(func() {
			mapper.CreateMap(func(user *User) (*Staff, error) {
				if user == nil {
					return nil, userNotFoundErr
				}

				return &Staff{
					Name: user.Name,
				}, nil
			})
		})

		It("can map and return error", func() {
			staff, err := mapper.Map[User, Staff](nil)

			Expect(staff).Should(BeNil())
			Expect(err).Should(MatchError(userNotFoundErr))
		})
	})
})
