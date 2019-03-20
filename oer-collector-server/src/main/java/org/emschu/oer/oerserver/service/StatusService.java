package org.emschu.oer.oerserver.service;

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

import org.emschu.oer.core.model.repository.*;
import org.emschu.oer.core.service.ChannelReadonlyService;
import org.emschu.oer.oerserver.api.StatusResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.context.annotation.PropertySource;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

/**
 * this service handles cached status response providing
 */
@Service
@PropertySource("classpath:application.properties")
public class StatusService {

    @Value("${oer.latest_api_version}")
    private String latestApiVersion;

    @Autowired
    private ChannelRepository channelRepository;

    @Autowired
    private TvShowRepository tvShowRepository;

    @Autowired
    private ProgramEntryRepository programEntryRepository;

    @Autowired
    private ImageLinkRepository imageLinkRepository;

    @Autowired
    private TagRepository tagRepository;

    @Autowired
    private ChannelReadonlyService channelReadonlyService;

    @Cacheable(value = "oer_data_status")
    public StatusResponse getStatusResponse() {
        // fill status response
        StatusResponse response = new StatusResponse();
        response.setChannelCount(channelRepository.count());
        response.setTvShowCount(tvShowRepository.count());
        response.setProgramEntryCount(programEntryRepository.count());
        response.setImageLinksCount(imageLinkRepository.count());
        response.setTagCount(tagRepository.count());

        response.setTvChannels(channelReadonlyService.getAllSenders());

        response.setVersion(latestApiVersion);
        // FIXME
//        response.setCurrentlyUpdating(updaterService.isUpdating());
        response.setCurrentlyUpdating(false);
        response.setServerDateTime(LocalDateTime.now());
        return response;
    }
}
