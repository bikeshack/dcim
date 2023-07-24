package postgres

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/bikeshack/dcim/pkg/components"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// PGInsertComponent inserts a component into the database and returns the uuid generated within the database
func InsertComponent(db *sqlx.DB, component *components.Component) (string, error) {

	// Insert the component into the database and return the uuid generated within the database
	query, args, err := sqlx.Named("INSERT INTO components (xname, class, arch, net_type, role, flag) VALUES (:xname, :class, :arch, :net_type, :role, :flag) RETURNING uid", component)
	if err != nil {
		log.Info("Error preparing named query:", err)
		return "", err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query) //only if postgres
	row := db.QueryRowx(query, args...)
	err = row.Scan(&component.Uid)
	if err != nil {
		// TODO: Figure out how to handle the difference between a user error and a server error
		log.Debug("Error executing query: "+query+"\n  -", err)
		log.Debug("Args: ", args)
		return "", err
	}
	return component.Uid.String(), err
}

func UpdateComponent(db *sqlx.DB, component *components.Component) error {
	result, err := db.NamedExec("UPDATE components SET (class, arch, net_type, role, flag) = (:Class, :Arch, :NetType, :Role, :Flag) WHERE id = :ID", component)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	switch affected {
	case 0:
		// No rows were updated.  This is an error.
		return sql.ErrNoRows
	case 1:
		// Success!
		return nil
	default:
		// This should never happen
		return errors.New("Unexpected number of rows updated: " + strconv.FormatInt(affected, 10))
	}
}

func DeleteComponent(db *sqlx.DB, id string) error {

	_, err := db.Exec("DELETE FROM components WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// PGGetComponent returns a single component from the database based on the id or xname
func GetComponent(db *sqlx.DB, id string) (*components.Component, error) {
	component := &components.Component{}
	// See if the id looks like a uuid
	_, err := uuid.Parse(id)
	if err == nil {
		// lack of an error means this is a uuid
		log.Debug("Treating id as a uuid:" + id)
		err = db.Get(component, "SELECT uid, xname, class, arch, net_type, role, flag FROM components WHERE id = $1", id)
		// this can throw all kinds of errors.  sql.ErrNoRows is an interesting one.  Bubble them all up.
		if err != nil {
			return nil, err
		}
		// Success!
		return component, nil
	}
	// Now that we know id isn't a uuid, we can assume it is an xname
	log.Debug("Treating id as an xname:" + id)
	err = db.Get(component, "SELECT uid, xname, class, arch, net_type, role, flag FROM components WHERE xname = $1", id)
	// this can throw all kinds of errors.  sql.ErrNoRows is an interesting one.  Bubble them all up.
	if err != nil {
		return nil, err
	}
	// Success!
	return component, nil
}
