/*************************************************************************
 * Copyright 2019 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

package main

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	RegexTimestampProcessor string = `regextimestamp`
)

var (
	ErrEmptyRegex = errors.New("Empty regular expression")
	ErrEmptyMatch = errors.New("Empty TS-Match-Name")
	ErrNoSubexps  = errors.New("Must specify at least one subexpression")
)

type RegexTimestampConfig struct {
	Regex                     string // the regular expression to apply to the data
	TS_Match_Name             string // the submatch which contains the timestamp
	Timestamp_Format_Override string
	Timezone_Override         string
	Assume_Local_Timezone     bool
}

func NewRegexTimestampProcessor(cfg RegexTimestampConfig) (*RegexTimestamp, error) {
	if len(cfg.Regex) == 0 {
		return nil, ErrEmptyRegex
	}
	if len(cfg.TS_Match_Name) == 0 {
		return nil, ErrEmptyMatch
	}

	re, err := regexp.Compile(cfg.Regex)
	if err != nil {
		return nil, err
	}
	subexps := re.SubexpNames()
	if len(subexps) == 0 {
		return nil, ErrNoSubexps
	}
	var ok bool
	var idx int
	for i, n := range subexps {
		if n == cfg.TS_Match_Name {
			idx = i
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("Specified TS-Match-Name=%v not found in regular expression: %v", cfg.TS_Match_Name, cfg.Regex)
	}

	tcfg := Config{
		FormatOverride: cfg.Timestamp_Format_Override,
	}
	tg, err := NewTimeGrinder(tcfg)
	if err != nil {
		return nil, err
	}
	if cfg.Assume_Local_Timezone {
		tg.SetLocalTime()
	}
	if cfg.Timezone_Override != `` {
		err = tg.SetTimezone(cfg.Timezone_Override)
		if err != nil {
			return nil, err
		}
	}

	return &RegexTimestamp{
		RegexTimestampConfig: cfg,
		re:                   re,
		matchidx:             idx,
		tg:                   tg,
	}, nil
}

type RegexTimestamp struct {
	RegexTimestampConfig
	re       *regexp.Regexp
	matchidx int
	tg       *TimeGrinder
}

func (rt *RegexTimestamp) Process(ent *Entry) (rset []*Entry, err error) {
	matches := rt.re.FindSubmatch(ent.Data)
	// grab the extraction
	if len(matches) < rt.matchidx {
		// no extraction, skip it
		rset = []*Entry{ent}
		return
	}
	extracted, ok, err := rt.tg.Extract(matches[rt.matchidx])
	if err != nil {
		return nil, err
	} else if ok {
		// successful parse
		ent.TS = FromStandard(extracted)
	}
	rset = []*Entry{ent}
	return
}
