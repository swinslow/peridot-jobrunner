// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package jobrunner

import (
	"time"

	"github.com/swinslow/peridot-db/pkg/datastore"
)

// Env is the environment for the web handlers.
type Env struct {
	Db datastore.Datastore
}

// SetupEnv sets up systems (such as the data store) and variables
// (such as the JWT signing key) that are used across web requests.
func SetupEnv() (*Env, error) {
	// set up datastore
	db, err := datastore.NewDB("host=db sslmode=disable dbname=dev user=postgres-dev")
	if err != nil {
		return nil, err
	}

	// don't init database tables if they don't yet exist;
	// that'll be the API's responsibility
	// FIXME is this the right way to handle this?
	// FIXME need to wait until tables are ready?
	// err = datastore.InitNewDB(db)
	// if err != nil {
	// 	return nil, err
	// }

	// FIXME THIS IS BAD: just waits for the DB to be ready,
	// FIXME and assumes it will be after x seconds
	time.Sleep(5 * time.Second)

	env := &Env{
		Db: db,
	}
	return env, nil
}
