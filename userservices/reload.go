package usersystem

type Reloader interface {
	//Reload reload user data
	Reload(string) error
}

//Reload reload user data with given reloaders
func Reload(id string, reloaders ...Reloader) error {
	var finalerr error
	for k := range reloaders {
		if reloaders[k] != nil {
			err := reloaders[k].Reload(id)
			if err != nil {
				finalerr = err
			}
		}
	}
	return finalerr
}
