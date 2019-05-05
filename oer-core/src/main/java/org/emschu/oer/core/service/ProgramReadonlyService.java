package org.emschu.oer.core.service;

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

import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.repository.ChannelRepository;
import org.emschu.oer.core.model.repository.ProgramEntryRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Service
public class ProgramReadonlyService {

    @Autowired
    private ProgramEntryRepository programEntryRepository;

    @Autowired
    private ChannelRepository channelRepository;

    public List<ProgramEntry> getCompleteDailyProgram(LocalDate day) {
        LocalDateTime startDate = LocalDateTime.of(day, LocalTime.of(0,0,0));
        LocalDateTime endDate = LocalDateTime.of(day, LocalTime.of(0,0,0))
                .plus(1, ChronoUnit.DAYS);
        return getProgram(startDate, endDate);
    }

    public List<ProgramEntry> getDayProgramByChannelId(LocalDate day, Long channelId) {
        LocalDateTime startDate = LocalDateTime.of(day, LocalTime.of(0,0,0));
        LocalDateTime endDate = LocalDateTime.of(day, LocalTime.of(0,0,0))
                .plus(1, ChronoUnit.DAYS);
        return getProgramOfChannel(startDate, endDate, channelId);
    }

    public List<ProgramEntry> getProgramOfChannel(LocalDateTime startDateTime, LocalDateTime endDateTime, Long channelId) {
        if (channelId == null) {
            throw new IllegalArgumentException("null channel given");
        }
        Optional<Channel> channelOptional = channelRepository.findById(channelId);
        if (!channelOptional.isPresent()) {
            // invalid channel
            return new ArrayList<>();
        }
        return programEntryRepository.findByStartDateTimeIsBetweenAndChannelOrderByStartDateTime(startDateTime, endDateTime, channelOptional.get());
    }

    public List<ProgramEntry> getProgram(LocalDateTime startDateTime, LocalDateTime endDateTime) {
        return programEntryRepository.findByStartDateTimeIsBetweenOrderByStartDateTime(startDateTime, endDateTime);
    }

    public Optional<ProgramEntry> getProgramEntry(Long entryId) {
        if (entryId == null) {
            throw new IllegalArgumentException("null entry id given");
        }
        return programEntryRepository.findById(entryId);
    }
}
