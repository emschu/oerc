package org.emschu.oer.collector.reader.parser.srf;

/*-
 * #%L
 * oer-collector-project
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

import org.emschu.oer.collector.reader.AbstractTvShowParser;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.parser.TvShowParserException;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.core.util.Hasher;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.logging.Logger;

@Component("srfTvShowParser")
public class TvShowParser extends AbstractTvShowParser {
    private static final String ALL_TV_SHOWS = "https://www.srf.ch/sendungen-a-z";
    private static final Logger LOG = Logger.getLogger(TvShowParser.class.getName());
    private static final String SRF_BASE_PAGE_URL = "https://www.srf.ch";

    @Override
    public List<TvShow> getEntries() throws TvShowParserException {
        List<TvShow> tvShows = new ArrayList<>();

        final Document body = Fetcher.fetchDocument(ALL_TV_SHOWS, "body");
        if (body == null) {
            LOG.warning(String.format("No content in url '%s'", ALL_TV_SHOWS));
            return new ArrayList<>();
        }
        final Elements select = body.select("#letters .row ul li");
        if (select == null) {
            throw new TvShowParserException("no entries found for url: " + ALL_TV_SHOWS);
        }

        for (Element e : select) {
            TvShow tvShow = new TvShow();
            final Elements titleElement = e.select("h3");

            if (titleElement.hasClass("radio-station")) {
                LOG.fine("Skipping srf radio station: " + titleElement.text());
                continue;
            }

            String title = titleElement.text();
            String url = e.select("h3 a").attr("href");
            String tvShowId = e.attr("data-show-teaser-id");
            String tvShowImage = e.attr("data-image-s3-retina");

            if (!url.contains("http")) {
                url = SRF_BASE_PAGE_URL + url;
            }

            tvShow.setTitle(title);
            tvShow.setHomePage(url);
            tvShow.setUrl(url);

            String tvShowHash = Hasher.getHash(tvShowId + ":" + getAdapterFamily());

            registerKey(tvShowId, tvShowHash);

            tvShow.setAdditionalId(tvShowId);
            tvShow.setTechnicalId(tvShowHash);
            tvShow.setAdapterFamily(getAdapterFamily());
            tvShow.setImageUrl(tvShowImage);

            tvShows.add(tvShow);
        }
        return tvShows;
    }

    @Override
    public void postProcessEntry(TvShow tvShow) throws TvShowParserException {
        // useful information is captured within the getElements() method above.
        // the angular content we get out of tvShow.url is not really parseable by jsoup.
    }

    @Override
    public Channel.AdapterFamily getAdapterFamily() {
        return Channel.AdapterFamily.SRF;
    }

    @Override
    public void cleanup() {
        // nothing to do here
    }
}
