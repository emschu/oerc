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

import org.emschu.oer.collector.reader.parser.ard.ARDReader;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.parser.orf.ORFReader;
import org.emschu.oer.collector.reader.parser.srf.SRFReader;
import org.emschu.oer.collector.reader.parser.zdf.ZDFReader;
import org.emschu.oer.core.service.EnvService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Caching;
import org.springframework.stereotype.Component;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.logging.Logger;

/**
 * central scheduling service implementation of this module
 */
@Component
public class UpdaterService {

    private static final Logger LOG = Logger.getLogger(UpdaterService.class.getName());

    private boolean isUpdating = false;
    private LocalDateTime updateStartDate;

    @Autowired
    private ARDReader ardReader;

    @Autowired
    private ZDFReader zdfReader;

    @Autowired
    private ORFReader orfReader;

    @Autowired
    private SRFReader srfReader;

    @Autowired
    private EnvService envService;

    @Value(value = "${oer.collector.skip_ard}")
    private boolean skipArd;

    @Value(value = "${oer.collector.skip_zdf}")
    private boolean skipZdf;

    @Value(value = "${oer.collector.skip_orf}")
    private boolean skipOrf;

    @Value(value = "${oer.collector.skip_srf}")
    private boolean skipSrf;

    @Value(value = "${oer.collector.proxy_host}")
    private String proxyHost;

    @Value(value = "${oer.collector.proxy_port}")
    private String proxyPort;

    public void fetchNewTVProgram() throws ParserException, InterruptedException {
        if (envService.isTestMode()) {
            LOG.info("NO data is updated in test profile. Skipping update.");
            return;
        }

        if (!connectionCheck()) {
            if (proxyHost != null || proxyPort != null) {
                LOG.warning(String.format("No internet connection available. Proxy seems not up at: %s:%s", proxyHost, proxyPort));
            } else {
                LOG.warning("No internet connection available.");
            }
            return;
        }

        startUpdating();
        LOG.info("Fetching new tv program - Start");
        try {
            ardReader.execute();
            zdfReader.execute();
            orfReader.execute();
            srfReader.execute();
        } catch (Exception e) {
            LOG.throwing(UpdaterService.class.getName(), "fetchNewTvProgram", e);
            throw e;
        }

        // log metrics
        LOG.info("Total request count: " + Fetcher.getCounter());
        if (this.updateStartDate == null) {
            LOG.warning("Illegal state: unknown update start date.");
        } else {
            LOG.info("Time needed for update: " + ChronoUnit.SECONDS.between(this.updateStartDate, LocalDateTime.now()) + " seconds");
        }
        LOG.info("Fetching new tv program - Finish");
        endUpdating();
    }

    /**
     * this method does an (initial) connection check to target websites
     *
     * @return
     */
    public boolean connectionCheck() {
        try {
            if (!skipArd) {
                Fetcher.fetchDocument("https://programm.ard.de", "body");
            }
            if (!skipZdf) {
                Fetcher.fetchDocument("https://www.zdf.de/live-tv", "body");
            }
            if (!skipOrf) {
                Fetcher.fetchDocument("https://orf.at/", "body");
            }
            if (!skipSrf) {
                Fetcher.fetchDocument("https://srf.ch/", "body");
            }
            return true;
        } catch (IllegalStateException ignored) { }
        return false;
    }

    /**
     * needs to be public for caching. do not use this method outside of the class.
     * NO EXTERNAL USE
     */
    @Caching(evict = {
            @CacheEvict("oer_data_status"),
    })
    public void startUpdating() {
        isUpdating = true;
        this.updateStartDate = LocalDateTime.now();
    }

    /**
     * needs to be public for caching. do not use this method outside of the class
     * NO EXTERNAL USE
     */
    @Caching(evict = {
            @CacheEvict("oer_data_status"),
    })
    public void endUpdating() {
        isUpdating = false;
        this.updateStartDate = null;
    }

    /**
     * @return nullable
     */
    public LocalDateTime getUpdateStartDate() {
        return updateStartDate;
    }

    /**
     * current updating state
     *
     * @return
     */
    public boolean isUpdating() {
        return isUpdating;
    }
}
