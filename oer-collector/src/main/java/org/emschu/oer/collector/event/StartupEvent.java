package org.emschu.oer.collector.event;

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

import org.emschu.oer.collector.reader.ZdfApiFetcher;
import org.emschu.oer.collector.service.CacheService;
import org.emschu.oer.collector.service.ChannelService;
import org.emschu.oer.collector.service.UpdaterService;
import org.emschu.oer.core.service.EnvService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Component;

import java.util.Arrays;
import java.util.logging.Logger;

@Component
public class StartupEvent implements ApplicationListener<ApplicationReadyEvent> {

    private static final Logger LOG = Logger.getLogger(StartupEvent.class.getName());

    @Autowired
    private ChannelService channelService;

    @Autowired
    private CacheService cacheService;

    @Autowired
    private EnvService envService;

    @Autowired
    private UpdaterService updaterService;

    @Value(value = "${oer.collector.proxy_host}")
    private String proxyHost;

    @Value(value = "${oer.collector.proxy_port}")
    private String proxyPort;

    @Override
    public void onApplicationEvent(ApplicationReadyEvent applicationReadyEvent) {
        LOG.info("OER-COLLECTOR is starting up and initializes channel families");
        LOG.info("Active profiles: " + Arrays.toString(envService.getActiveProfiles()));

        if (envService.isTestMode()) {
            LOG.info("Collector is in TEST mode!");
        }

        if (updaterService.connectionCheck()) {
            ZdfApiFetcher.setCurrentZdfApiKey(cacheService.getZdfApiKey());

            channelService.init();
        } else {
            if (proxyHost != null || proxyPort != null) {
                LOG.warning(String.format("No internet connection available. Proxy seems not up at: %s:%s", proxyHost, proxyPort));
            } else {
                LOG.warning("No internet connection available.");
            }
        }
        // .. in several secs the update service is starting
        /* @see org.emschu.oer.collector.service.UpdaterService */
    }
}
