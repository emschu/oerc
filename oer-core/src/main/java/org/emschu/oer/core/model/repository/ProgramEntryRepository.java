package org.emschu.oer.core.model.repository;

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
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Repository
public interface ProgramEntryRepository extends CrudRepository<ProgramEntry, Long> {
    public boolean existsByTechnicalId(String technicalId);
    public ProgramEntry getByTechnicalIdAndChannel(String technicalId, Channel channel);
    public Optional<ProgramEntry> findByTechnicalIdAndAdapterFamily(String technicalId, Channel.AdapterFamily adapterFamily);
    public List<ProgramEntry> findByStartDateTimeAfterAndEndDateTimeBefore(LocalDateTime startDateTime, LocalDateTime endDateTime);
    public List<ProgramEntry> findByStartDateTimeAfterAndEndDateTimeBeforeAndChannel(LocalDateTime startDateTime, LocalDateTime endDateTime, Channel channel);
    public Iterable<ProgramEntry> getAllByTechnicalIdIsInAndAdapterFamily(List<String> technicalIdList, Channel.AdapterFamily adapterFamily);
}
