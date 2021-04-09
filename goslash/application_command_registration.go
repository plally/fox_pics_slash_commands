package goslash

import (
	log "github.com/sirupsen/logrus"
)

// ensure the global command list is up to date with what is listed by this app internally
func (app *Application) SyncGlobal() {
	commands, err := app.Session.ApplicationCommands(app.ClientID, "")
	if err != nil {
		log.WithError(err).Warn("error syncing global commands")
	}
	for _, command := range commands {
		localCommand := app.GetCommand(command.Name)
		if localCommand == nil {
			app.DeleteGlobal(command.ID)
			continue
		}
	}
	// TODO check command equality to determine which commands should be registered
	app.RegisterAllGlobal()
}

func (app *Application) RegisterGlobal(command *Command) (*Command, error) {
	newCommand, err := app.Session.ApplicationCommandCreate(app.ClientID, "", &command.ApplicationCommand)
	if err != nil {
		log.WithField("error", err).Info("Error occurred while registering global command")
		return command, err
	}

	app.Commands[command.Name] = command
	command.isGlobal = true
	command.ApplicationCommand = *newCommand

	return command, err
}

func (app *Application) RegisterGuild(guildid string, command *Command) (*Command, error) {
	_, err := app.Session.ApplicationCommandCreate(app.ClientID, guildid, &command.ApplicationCommand)
	if err != nil {
		return nil, err
	}

	return command, err
}

func (app *Application) RegisterAllGuild(guildid string) error {
	for _, command := range app.Commands {
		_, err := app.RegisterGuild(guildid, command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) RegisterAllGlobal() error {
	for _, command := range app.Commands {
		_, err := app.RegisterGlobal(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) DeleteGlobal(commandId string) error {
	return app.Session.ApplicationCommandDelete(app.ClientID, "", commandId)
}
