package org.emschu.oer.oerserver.controller;

/*-
 * #%L
 * oer-server
 * %%
 * Copyright (C) 2019 emschu[aet]mailbox.org
 * %%
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * #L%
 */

import io.swagger.annotations.Api;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.service.ChannelReadonlyService;
import org.emschu.oer.core.service.ProgramReadonlyService;
import org.emschu.oer.oerserver.api.StatusResponse;
import org.emschu.oer.oerserver.service.StatusService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.logging.Logger;

@Api
@RestController
@RequestMapping(path = "v1")
public class StatusController {

    @Autowired
    private StatusService statusService;

    @Autowired
    private ChannelReadonlyService channelReadonlyService;

    @Autowired
    private ProgramReadonlyService programService;

    private static final Logger LOG = Logger.getLogger(StatusController.class.getName());

    // status action
    @GetMapping("/status")
    public StatusResponse statusAction() {
        return statusService.getStatusResponse();
    }

    // channel actions
    @GetMapping("/channels")
    public List<Channel> channelListAction() {
        return channelReadonlyService.getAllSenders();
    }

    @GetMapping("/channel/{id}")
    public Channel channelSingleAction(@PathVariable Long id) {
        return channelReadonlyService.getChannel(id);
    }

    // daily program actions
    @GetMapping("/program/daily")
    public ResponseEntity<List<ProgramEntry>> allDailyProgramEntries() {
        return new ResponseEntity<>(programService.getCompleteDailyProgram(LocalDate.now()), HttpStatus.ACCEPTED);
    }

    @GetMapping("/program/daily/{channelId}")
    public ResponseEntity<List<ProgramEntry>> dailyProgramEntriesByChannel(@PathVariable Long channelId) {
        if (channelId == null) {
            return ResponseEntity.badRequest().build();
        }
        return new ResponseEntity<>(programService.getDayProgramByChannelId(LocalDate.now(), channelId), HttpStatus.ACCEPTED);
    }

    @GetMapping("/program")
    public ResponseEntity<List<ProgramEntry>> programEntries(@RequestParam(name = "from", required = true) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime fromDate,
                                                             @RequestParam(name = "to", required = true) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime toDate,
                                                             @RequestParam(name = "channelId", required = false) Long channelId) {
        if (channelId == null) {
            return ResponseEntity.ok(programService.getProgram(fromDate, toDate));
        }
        return ResponseEntity.ok(programService.getProgramOfChannel(fromDate, toDate, channelId));
    }

    @GetMapping("/program/entry/{entryId}")
    public ResponseEntity<ProgramEntry> programEntry(@PathVariable Long entryId) {
        if (entryId == null) {
            return ResponseEntity.badRequest().build();
        }
        Optional<ProgramEntry> programEntryOptional = programService.getProgramEntry(entryId);
        if (programEntryOptional.isPresent()) {
            return ResponseEntity.ok(programEntryOptional.get());
        }
        return ResponseEntity.notFound().build();
    }
}
