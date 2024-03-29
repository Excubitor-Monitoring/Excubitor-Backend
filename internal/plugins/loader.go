package plugins

import (
	"fmt"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/config"
	ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/logging"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared"
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
	"strings"
)

var logger logging.Logger
var loadablePlugins []string

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MODULE_PLUGIN",
	MagicCookieValue: "EXCUBITOR",
}

func LoadPlugins() error {
	logger = logging.GetLogger()
	logger.Debug("Loading plugins...")

	pluginFolder, err := os.ReadDir("plugins")
	if err != nil {
		logger.Error("Error occurred upon trying to read plugins.")
		return fmt.Errorf("loading plugins: %w", err)
	}

	for _, pluginEntry := range pluginFolder {
		if pluginEntry.IsDir() {
			items, err := os.ReadDir("plugins/" + pluginEntry.Name())
			if err != nil {
				logger.Error(fmt.Sprintf("Error on loading plugin %s. Skipping... Reason: %s", pluginEntry.Name(), err))
				continue
			}

			for _, item := range items {
				if strings.HasSuffix(item.Name(), ".plugin") {
					loadablePlugins = append(loadablePlugins, "plugins/"+pluginEntry.Name()+"/"+item.Name())
					logger.Debug(fmt.Sprintf("Added plugin %s to loadable plugins.", item.Name()))
				}
			}
		} else {
			logger.Info(fmt.Sprintf("Unknown item %s in plugins folder.", pluginEntry.Name()))
		}
	}

	return nil
}

func InitPlugins() error {
	for _, pl := range loadablePlugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins: map[string]plugin.Plugin{
				"module": &shared.ModulePlugin{},
			},
			Cmd: exec.Command("./" + pl),
			Logger: NewLogWrapper(
				logger,
				strings.Split(pl, "/")[1],
				logging.GetLogLevelByString(config.GetConfig().String("logging.log_level")),
			),
		})

		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		rawPlugin, err := rpcClient.Dispense("module")
		if err != nil {
			return err
		}

		loadedPlugin := rawPlugin.(shared.ModuleProvider)

		module := modules.NewModule(
			loadedPlugin.GetName(),
			loadedPlugin.GetVersion(),
			nil,
			func() {
				messages := loadedPlugin.TickFunction()
				for _, msg := range messages {
					ctx.GetContext().GetBroker().Publish(msg.Monitor, msg.Body)
				}
			},
		)

		for _, component := range loadedPlugin.GetComponents() {
			module.Components = append(module.Components,
				modules.Component{
					TabName: component.TabName,
					JSFile:  "static/external/" + loadedPlugin.GetName() + "/" + component.JSFile,
					Tag:     component.Tag,
				})
		}

		frontendGetters[module.Name] = loadedPlugin.GetComponentFile

		ctx.GetContext().RegisterModule(module)
	}

	return nil
}
