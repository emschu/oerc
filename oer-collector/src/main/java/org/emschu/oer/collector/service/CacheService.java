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

import org.emschu.oer.collector.reader.parser.zdf.ProgramEntryParser;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.logging.Logger;

/**
 * this class encapsulates methods/functionality which should be cached
 */
@Service
public class CacheService {

    private static final Logger LOG = Logger.getLogger(CacheService.class.getName());

    @Value("${oer.collector.skip_zdf}")
    private boolean skipZdf;

    @Cacheable("oer_data_api_key")
    public String getZdfApiKey() {
        LOG.info("Starting cached method of retrieving zdf api key");

        if (skipZdf) {
            LOG.info("skip_zdf = true, will not connect to retrieve api key: SKIPPING");
            return "";
        }
        ProgramEntryParser.ZDFScraper zdfScraper = new ProgramEntryParser.ZDFScraper();
        String apiKey = zdfScraper.retrieveApiKey();
        LOG.info("retrieved zdf api key: " + apiKey);
        return apiKey;
    }
}
