// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"strings"

	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/dankgrinder/dankgrinder/instance/scheduler"
)

func (in *Instance) gift(msg discord.Message) {
	trigger := in.sdlr.AwaitResumeTrigger()
	if trigger == nil || !strings.Contains(trigger.Value, shopBaseCmdValue) {
		return
	}
	if in == in.Master {
		in.sdlr.Resume()
		return
	}
	if !exp.gift.Match([]byte(msg.Embeds[0].Title)) || !exp.shop.Match([]byte(trigger.Value)) {
		in.sdlr.Resume()
		return
	}
	amount := strings.Replace(exp.gift.FindStringSubmatch(msg.Embeds[0].Title)[1], ",", "", -1)
	item := exp.shop.FindStringSubmatch(trigger.Value)[1]

	in.iteratedItems++

	if amount == "0" {
		in.sdlr.Resume()
		return
	}

	// If not the last item to send in gift item list or does not exceed 
	// max items per trade append the items to the list of items to send
	// add one as it does not count the current item 
	if in.currentTradeItems+1 < in.Features.MaxItemsPerTrade || in.iteratedItems == in.totalTradeItems {
		in.tradeList += tradeItemListValue(amount, item)
		in.currentTradeItems++
		in.sdlr.Resume()
		return
	} else {
		in.tradeList += tradeItemListValue(amount, item)
	}

	// ResumeWithCommandOrPrioritySchedule is not necessary in this case because
	// the scheduler has to be awaiting resume. AwaitResumeTrigger returns "" if
	// the scheduler isn't awaiting resume which causes this function to return.
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Value: tradeCmdValue(in.tradeList, in.Master.Client.User.ID),
		Log:   "gifting items - starting trade",
		AwaitResume: true,
	})

	in.tradeList = ""
	in.currentTradeItems = 0
}

func (in *Instance) confirmTrade(msg discord.Message) {
	in.sdlr.ResumeWithCommand(&scheduler.Command{
		Actionrow: 1,
		Button: 2,
		Message: msg,
		Log: "gifting items - accepting trade as sender",
		AwaitResume: true,
	})
}

func (in *Instance) confirmTradeAsMaster(msg discord.Message) {
	// If trade request mentioning master is sent, priority schedule a click on accept
	in.Master.sdlr.PrioritySchedule(&scheduler.Command{
		Actionrow: 1,
		Button: 2,
		Message: msg,
		Log: "gifting items - accepting trade as master",
	})
}