/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package server

import (
	"fmt"
	"net/http"

	"github.com/SuperGreenLab/EmergencyUpgrader/internal/server/routes"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("serverport", "8081", "Server port")
)

func init() {
	viper.SetDefault("ServerPort", "8081")
}

// Start starts the server
func Start() {
	router := httprouter.New()

	routes.Init(router)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("ServerPort")), cors.AllowAll().Handler(router)))
	}()
}
