package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	createUsage string = "Usage: migrator create [NAME]\n" +
		"The create command is used to create two timestamped migration files, " +
		"one for the up migration and one for the down migration\n"
)

// TODO: Check file name and write sql file accordingly:
//
// create_xxx => CREATE TABLE xxx ();
// update_xxx | change_xxx => ALTER TABLE xxx ();
// drop_xxx | delete_xxx => DROP TABLE xxx;
func NewCreateCommand() *Command {
	c := NewCommand(
		"create",
		createUsage,
		func(args []string) error {
			set := flag.NewFlagSet("create", flag.ContinueOnError)
			dirPtr := set.String("d", "database", "Path to the migration directory")
			err := set.Parse(args)
			if err != nil {
				if errors.Is(err, flag.ErrHelp) {
					return nil
				}
				return err
			}
			if set.NArg() != 1 {
				return errors.New("need exactly one argument")
			}
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			dest := path.Join(wd, *dirPtr)
			err = os.MkdirAll(dest, 0777)
			if err != nil {
				if !errors.Is(err, fs.ErrExist) {
					return err
				}
			}
			name := set.Arg(0)
			upPath, downPath := generateUpAndDownFileNames(strconv.Itoa(int(time.Now().Unix())), name)
			upPath = path.Join(dest, upPath)
			downPath = path.Join(dest, downPath)
			//OpenFile instead of Create makes sure that no existing files are truncated.
			//This is only a theoretical case as in practice, each file has a timestamp prefix and is thus not likely
			//to be truncated
			f, err := os.OpenFile(upPath, os.O_CREATE|os.O_RDONLY, 0666)
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
			f, err = os.OpenFile(downPath, os.O_CREATE|os.O_RDONLY, 0666)
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
			return nil
		},
	)
	return c
}

func generateUpAndDownFileNames(prefix, name string) (string, string) {
	return generateGeneralMigrationFileName(prefix, name, "up"), generateGeneralMigrationFileName(prefix, name, "down")
}

func generateGeneralMigrationFileName(prefix, name, suffix string) string {
	return fmt.Sprintf("%s_%s.%s.sql", prefix, name, suffix)
}
