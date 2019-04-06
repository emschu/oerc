package org.emschu.oer.collector.reader.parser.orf;

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
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;

@Component("orfTvShowParser")
public class TvShowParser extends AbstractTvShowParser {

    @Override
    public List<TvShow> getEntries() throws TvShowParserException {
        final Document body = Fetcher.fetchDocument("https://tvthek.orf.at/profiles", "body");
        final Elements select = body.select("article.b-teaser");

        List<TvShow> tvShowList = new ArrayList<>();
        for (Element singleTeaser : select) {
            final String tvShowTitle = singleTeaser.select("a").attr("title");
            final String tvShowUrl = singleTeaser.select("a").attr("href");

            handleTvShow(tvShowList, tvShowTitle, tvShowUrl, null);
        }

        return tvShowList;
    }

    @Override
    public void postProcessEntry(TvShow tvShow) throws TvShowParserException {
        // there is almost no possibility to get more useful information atm from orf
    }

    @Override
    public Channel.AdapterFamily getAdapterFamily() {
        return Channel.AdapterFamily.ORF;
    }

    @Override
    public void cleanup() {
        this.clear();
    }
}
