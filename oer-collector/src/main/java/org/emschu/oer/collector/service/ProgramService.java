package org.emschu.oer.collector.service;

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
import org.emschu.oer.core.model.repository.ProgramEntryRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Optional;

@Service
public class ProgramService {

    @Autowired
    private ProgramEntryRepository programEntryRepository;

    public Optional<ProgramEntry> findProgramEntryByTechnicalId(String technicalId, Channel.AdapterFamily adapterFamily) {
        return programEntryRepository.findByTechnicalIdAndAdapterFamily(technicalId, adapterFamily);
    }

    public Optional<ProgramEntry> nearestProgramEntryInFuture(LocalDateTime startDateTime, Channel channel) {
        return programEntryRepository.findFirstByChannelAndStartDateTimeIsAfterOrderByStartDateTime(channel, startDateTime);
    }
}
