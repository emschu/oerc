package org.emschu.oer.collector.reader.parser.ard;

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

import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.emschu.oer.collector.reader.AbstractTvShowParser;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.parser.TvShowParserException;
import org.emschu.oer.collector.reader.parser.TvShowParserInterface;
import org.emschu.oer.collector.service.ProgramService;
import org.emschu.oer.collector.util.EidExtractor;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.core.util.Hasher;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.logging.Logger;

@Component(value = "ardTvShowParser")
public class TvShowParser extends AbstractTvShowParser {

    private static final String INTER_URL_PART = "Sendungen-von-A-bis-Z";
    private static final String ALL_CHANNEL_URL = "https://programm.ard.de/TV/" + INTER_URL_PART + "/Startseite?page=&char=all";

    @Autowired
    private ProgramService programService;

    private static final Logger LOG = Logger.getLogger(TvShowParser.class.getName());

    @Override
    public List<TvShow> getEntries() throws TvShowParserException {
        List<TvShow> tvShowList = new ArrayList<>();
        Document jsoupDoc = Fetcher.fetchDocument(ALL_CHANNEL_URL, ".az .con div.box");
        Elements selection = jsoupDoc.select(".az .con div.box");
        if (selection.isEmpty()) {
            throw new TvShowParserException("no tv shows found for url: " + ALL_CHANNEL_URL);
        }

        for (Element tvShowElement : selection) {
            String title = tvShowElement.select("img").attr("title").trim();
            String url = tvShowElement.select("a").attr("href").trim();

            if (!url.contains(INTER_URL_PART)) {
                // 'whitelist' entries by url
                continue;
            }
            handleTvShow(tvShowList, title, url, ProgramEntryParser.ARD_HOST);
        }
        return tvShowList;
    }

    @Override
    public void postProcessEntry(TvShow tvShowEntry) throws TvShowParserException {
        String url = tvShowEntry.getUrl();
        if (url == null || url.isEmpty()) {
            throw new TvShowParserException("no url in tv show found");
        }
        LOG.info("Fetching: " + url);
        final String mainTvShowImageElement = ".summary .gallery img";
        final Document document = Fetcher.fetchDocument(url, mainTvShowImageElement);

        String imageUrl = ProgramEntryParser.ARD_HOST + document.select(mainTvShowImageElement).attr("src");
        tvShowEntry.setImageUrl(imageUrl);

        Elements realUrl = document.select(".gal-link").select("a");
        if (!realUrl.isEmpty()) {
            if (realUrl.text().equals("Sendungshomepage")) {
                tvShowEntry.setHomePage(realUrl.attr("href"));
            }
        }

        linkRelatedProgramEntries(document, tvShowEntry);
    }

    protected void linkRelatedProgramEntries(Document document, TvShow tvShow) {
        int failureCounter = 0;
        Elements programEntryLinks = document.select(".bucket-row ul li[class^=eid]");
        if (programEntryLinks == null) {
            LOG.fine("No program entry links found");
            return;
        }
        for (Element singleProgramEntry : programEntryLinks) {
            String scrapedTechnicalId = singleProgramEntry.attr("class");
            scrapedTechnicalId = EidExtractor.extractEid(scrapedTechnicalId);
            if (scrapedTechnicalId.isEmpty()) {
                continue;
            }
            // if the program entry is known - we store it
            Optional<ProgramEntry> possibleProgramEntry =
                    programService.findProgramEntryByTechnicalId(scrapedTechnicalId, tvShow.getAdapterFamily());
            if (possibleProgramEntry.isPresent()) {
                if (tvShow.getRelatedProgramEntries().contains(possibleProgramEntry.get())) {
                    // avoid duplicates in db record and with collecting list
                    continue;
                }
                // failure counter is reset if there are entries
                failureCounter = 0;
                tvShow.getRelatedProgramEntries().add(possibleProgramEntry.get());
            } else {
                failureCounter++;
            }
            if (failureCounter > 8) {
                LOG.finest("Related program entry linkage skipped after 8 failures");
                // leave loop after 8 failures to speed up everything
                break;
            }
        }
        int linkedTagSize = 0;
        if (tvShow.getRelatedProgramEntries() != null) {
            linkedTagSize = tvShow.getRelatedProgramEntries().size();
        }
        LOG.info(String.format("Detected %d program entries, linked %d", programEntryLinks.size(), linkedTagSize));
    }

    @Override
    public Channel.AdapterFamily getAdapterFamily() {
        return Channel.AdapterFamily.ARD;
    }

    @Override
    public void cleanup() {
        this.clear();
    }
}
