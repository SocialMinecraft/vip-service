package database

type Account struct {
	Id            int
	UserId        string
	MinecraftUuid string
}

func (r *Db) AddAccount(account Account) error {
	_, err := r.db.Query(
		"INSERT INTO accounts (user_id, minecraft_uuid) VALUES ($1, $2) ON CONFLICT (minecraft_uuid) DO NOTHING",
		account.UserId,
		account.MinecraftUuid,
	)
	return err
}

func (r *Db) RemoveAccountByUuid(uuid string) error {
	_, err := r.db.Query(
		"DELETE FROM accounts WHERE minecraft_uuid = $1",
		uuid,
	)
	return err
}
