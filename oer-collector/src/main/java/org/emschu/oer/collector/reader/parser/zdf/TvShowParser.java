package org.emschu.oer.collector.reader.parser.zdf;

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
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.Tag;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.ZdfApiFetcher;
import org.emschu.oer.collector.reader.parser.TvShowParserException;
import org.emschu.oer.collector.service.TagService;
import org.emschu.oer.core.util.Hasher;
import org.emschu.oer.zdf_api.model.EditorialTag;
import org.emschu.oer.zdf_api.model.ZdfTvShowResponseModel;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.List;
import java.util.logging.Logger;

@Component(value = "zdfTvShowParser")
public class TvShowParser extends AbstractTvShowParser {

    private static final String ZDF_TV_SHOW_LINK = "https://www.zdf.de/sendungen-a-z?group=";
    private static String[] tvShowGroups = {
            "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
            "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
            "0+-+9"
    };
    private static final Logger LOG = Logger.getLogger(TvShowParser.class.getName());

    @Autowired
    private TagService tagService;

    @Override
    public List<TvShow> getEntries() throws TvShowParserException {
        List<TvShow> tvShowList = new ArrayList<>();

        collectEntries(tvShowList);

        return tvShowList;
    }

    private void collectEntries(List<TvShow> tvShowList) {
        for (int z = 0; z < tvShowGroups.length; z++) {
            String apiUrl = ZDF_TV_SHOW_LINK + tvShowGroups[z];
            Document zdfTvShowPage = Fetcher.fetchDocument(apiUrl, "article");
            Elements articleElements = zdfTvShowPage.select("article.b-content-teaser-item");
            int counter = 0;
            for (Element article : articleElements) {
                final Elements bPlusBtnElement = article.select(".b-plus-button");
                String tvShowTitle = bPlusBtnElement.attr("data-plusbar-title");
                String tvShowUrl = bPlusBtnElement.attr("data-plusbar-url");
                String tvShowApiPath = bPlusBtnElement.attr("data-plusbar-path"); // the zdf api path!
                String tvShowExternalId = bPlusBtnElement.attr("data-plusbar-external-id");
                String tvShowId = bPlusBtnElement.attr("data-plusbar-id");

                if (tvShowExternalId == null || tvShowExternalId.isEmpty()) {
                    throw new IllegalStateException("no externalId found");
                }
                if (tvShowId == null || tvShowId.isEmpty()) {
                    throw new IllegalStateException("no tv show id found");
                }

                // skip duplicates
                final String tvShowHash = Hasher.getHash(tvShowExternalId + tvShowId);
                if (isKeyRegistered(tvShowId)) {
                    LOG.info(String.format("Detected duplicate tv show: '%s'", tvShowTitle));
                    continue;
                }

                TvShow tvShow = new TvShow();
                tvShow.setTitle(tvShowTitle);
                tvShow.setHomePage(tvShowUrl);
                tvShow.setUrl(tvShowApiPath);

                // "register" this id
                registerKey(tvShowId, tvShowHash);
                tvShow.setTechnicalId(tvShowHash);
                tvShow.setAdapterFamily(getAdapterFamily());

                tvShowList.add(tvShow);
                counter++;
            }
            LOG.info(String.format("Detected %d tv shows of zdf page '%s'", counter, apiUrl));
        }
    }

    @Override
    public void postProcessEntry(TvShow tvShow) throws TvShowParserException {
        String url = tvShow.getUrl();
        if (url == null || url.isEmpty()) {
            throw new TvShowParserException("no url in tv show found");
        }
        final String urlToFetch = ProgramEntryParser.ZDFScraper.ZDF_API_HOST + "/content/documents/zdf" + tvShow.getUrl();
        ZdfTvShowResponseModel tvShowResponseModel =
                ZdfApiFetcher.getSingleTvShow(urlToFetch);

        // filter false/null responses
        if (tvShowResponseModel == null
                || tvShowResponseModel.getTeaserImageRef() == null
                || tvShowResponseModel.getTeaserImageRef().getLayouts() == null) {
            LOG.warning("Invalid zdf tv show response model given. perhaps invalid response or 403?");
            return;
        }

        tvShow.setImageUrl(tvShowResponseModel.getTeaserImageRef().getLayouts().getOriginal());

        List<Tag> tagListOfShow = new ArrayList<>();
        for (EditorialTag editorialTag : tvShowResponseModel.getEditorialTag()) {
            Tag tag = tagService.getOrCreateTag(editorialTag.getEditorialTagLabel());
            tagListOfShow.add(tag);
        }
        tvShow.setTags(tagListOfShow);

        String brandId = tvShowResponseModel.getHttpZdfDeRelsContentConfSection().getObjectId();
        tvShow.setAdditionalId(brandId);
    }

    @Override
    public Channel.AdapterFamily getAdapterFamily() {
        return Channel.AdapterFamily.ZDF;
    }

    @Override
    public void cleanup() {
        this.clear();
    }
}
