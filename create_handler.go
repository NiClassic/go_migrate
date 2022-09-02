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
			_, err = os.Create(upPath)
			if err != nil {
				return err
			}
			_, err = os.Create(downPath)
			if err != nil {
				return err
			}
			fmt.Printf("Created %s and %s\n", upPath, downPath)
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
