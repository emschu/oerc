package org.emschu.oer.collector.reader;

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

import org.emschu.oer.collector.reader.parser.TvShowParserInterface;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.core.util.Hasher;

import javax.validation.constraints.Null;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.logging.Logger;

public abstract class AbstractTvShowParser implements TvShowParserInterface {

    private static final Logger LOG = Logger.getLogger(AbstractTvShowParser.class.getName());
    private ConcurrentHashMap<String, String> usedTvShowIds = new ConcurrentHashMap<>();

    protected void registerKey(String tvShowId, String tvShowHash) {
        if (isKeyRegistered(tvShowId)) {
            throw new IllegalStateException(String.format("key '%s' is already registered", tvShowId));
        }
        usedTvShowIds.put(tvShowId, tvShowHash);
    }

    protected boolean isKeyRegistered(String key) {
        return usedTvShowIds.containsKey(key);
    }

    /**
     * method to create a tv show with its title + url
     *
     * @param tvShowList
     * @param title
     * @param url
     * @param urlPrefix
     */
    protected void handleTvShow(List<TvShow> tvShowList, String title, String url, @Null String urlPrefix) {
        // skip duplicates
        final String tvShowHash = Hasher.getHash(title + url);
        if (isKeyRegistered(url)) {
            LOG.info(String.format("Detected duplicate tv show: '%s'", title));
            return;
        }

        TvShow tvShow = new TvShow();
        tvShow.setTitle(title);
        if (urlPrefix != null) {
            tvShow.setUrl(urlPrefix + url);
        } else {
            tvShow.setUrl(url);
        }
        tvShow.setTechnicalId(tvShowHash);
        tvShow.setAdapterFamily(getAdapterFamily());
        if (!isKeyRegistered(tvShow.getUrl())) {
            registerKey(tvShow.getUrl(), tvShowHash);
        } else {
            LOG.warning("tv show is already registered: " + tvShow.getTitle() + " with url: " + tvShow.getUrl());
        }
        tvShowList.add(tvShow);
    }

    public void clear() {
        this.usedTvShowIds.clear();
    }
}
