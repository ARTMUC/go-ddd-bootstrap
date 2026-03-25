package userpersistence

import domainuser "starter/internal/domain/user"

// Translator maps between the User domain entity and the GORM persistence model.
type Translator struct{}

// ToModel converts a domain User to the GORM persistence model.
func (t Translator) ToModel(u *domainuser.User) *UserModel {
	return &UserModel{
		ID:   u.ID(),
		Name: u.Name(),
	}
}

// ToDomain converts a GORM persistence model back to the domain User entity.
func (t Translator) ToDomain(m *UserModel) *domainuser.User {
	return domainuser.Reconstitute(m.ID, m.Name)
}
