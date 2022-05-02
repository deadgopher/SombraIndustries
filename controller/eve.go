package controller

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type eve struct {
	base            string
	callback        string
	scope           string
	state           string
	client          string
	postDestination string
	creds           string
}

func (x *eve) registerRoutes(r *gin.RouterGroup) {
	r.GET("/", x.handleCallback)
}

func newEveController(r root) *eve {
	x := &eve{}
	x.base = "https://login.eveonline.com/v2/oauth/authorize/?response_type=code&"
	x.callback = "redirect_uri=" + url.QueryEscape("https://localhost:8080/api/eve/")
	x.scope = url.QueryEscape("scope=publicData esi-calendar.respond_calendar_events.v1 esi-calendar.read_calendar_events.v1 esi-location.read_location.v1 esi-location.read_ship_type.v1 esi-mail.organize_mail.v1 esi-mail.read_mail.v1 esi-mail.send_mail.v1 esi-skills.read_skills.v1 esi-skills.read_skillqueue.v1 esi-wallet.read_character_wallet.v1 esi-wallet.read_corporation_wallet.v1 esi-search.search_structures.v1 esi-clones.read_clones.v1 esi-characters.read_contacts.v1 esi-universe.read_structures.v1 esi-bookmarks.read_character_bookmarks.v1 esi-killmails.read_killmails.v1 esi-corporations.read_corporation_membership.v1 esi-assets.read_assets.v1 esi-planets.manage_planets.v1 esi-fleets.read_fleet.v1 esi-fleets.write_fleet.v1 esi-ui.open_window.v1 esi-ui.write_waypoint.v1 esi-characters.write_contacts.v1 esi-fittings.read_fittings.v1 esi-fittings.write_fittings.v1 esi-markets.structure_markets.v1 esi-corporations.read_structures.v1 esi-characters.read_loyalty.v1 esi-characters.read_opportunities.v1 esi-characters.read_chat_channels.v1 esi-characters.read_medals.v1 esi-characters.read_standings.v1 esi-characters.read_agents_research.v1 esi-industry.read_character_jobs.v1 esi-markets.read_character_orders.v1 esi-characters.read_blueprints.v1 esi-characters.read_corporation_roles.v1 esi-location.read_online.v1 esi-contracts.read_character_contracts.v1 esi-clones.read_implants.v1 esi-characters.read_fatigue.v1 esi-killmails.read_corporation_killmails.v1 esi-corporations.track_members.v1 esi-wallet.read_corporation_wallets.v1 esi-characters.read_notifications.v1 esi-corporations.read_divisions.v1 esi-corporations.read_contacts.v1 esi-assets.read_corporation_assets.v1 esi-corporations.read_titles.v1 esi-corporations.read_blueprints.v1 esi-bookmarks.read_corporation_bookmarks.v1 esi-contracts.read_corporation_contracts.v1 esi-corporations.read_standings.v1 esi-corporations.read_starbases.v1 esi-industry.read_corporation_jobs.v1 esi-markets.read_corporation_orders.v1 esi-corporations.read_container_logs.v1 esi-industry.read_character_mining.v1 esi-industry.read_corporation_mining.v1 esi-planets.read_customs_offices.v1 esi-corporations.read_facilities.v1 esi-corporations.read_medals.v1 esi-characters.read_titles.v1 esi-alliances.read_contacts.v1 esi-characters.read_fw_stats.v1 esi-corporations.read_fw_stats.v1 esi-characterstats.read.v1")
	x.state = "state=superSecretEveApiState"
	x.client = fmt.Sprintf("client_id=%v", os.Getenv("EVE_CLIENT_ID"))
	x.postDestination = "https://login.eveonline.com/v2/oauth/token"
	return x
}
func (x *eve) Link() string {
	return fmt.Sprintf(x.base+"%v&%v&%v&%v", x.callback, x.client, x.scope, x.state)
}

func (x *eve) handleCallback(c *gin.Context) {
	x.creds = "grant_type=authorization_code&code=" + c.Query("code")
	data := fmt.Sprintf("%v:%v", os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_SECRET_KEY"))
	encData := b64.URLEncoding.EncodeToString([]byte(data))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", x.postDestination, strings.NewReader(x.creds))
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	req.Header.Set("Authorization", "Basic "+encData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "login.eveonline.com")

	res, err2 := client.Do(req)
	if err2 != nil {
		respond(c, 500, false, foo(err))
		return
	}
	defer res.Body.Close()

	body, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		respond(c, 500, false, foo(err))
		return
	}

	respond(c, 200, true, string(body))
}
