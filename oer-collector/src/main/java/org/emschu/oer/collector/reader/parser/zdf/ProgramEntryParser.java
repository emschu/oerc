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

import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.ZdfApiFetcher;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.emschu.oer.core.model.ImageLink;
import org.emschu.oer.core.model.Tag;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.service.CacheService;
import org.emschu.oer.collector.service.ImageLinkService;
import org.emschu.oer.collector.service.TagService;
import org.emschu.oer.collector.service.TvShowService;
import org.emschu.oer.zdf_api.model.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.util.ConcurrentReferenceHashMap;

import java.io.*;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.time.*;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * a concrete implementation of {@link ProgramEntryParserInterface} to fetch and enrich ard program data
 */
@Component(value = "zdfProgramEntryParser")
public class ProgramEntryParser implements ProgramEntryParserInterface<ZdfBroadcast> {

    private static final String ZDF_HOST = "https://www.zdf.de";
    private static final Logger LOG = Logger.getLogger(ProgramEntryParser.class.getName());

    private Map<String, ZdfBroadcast> apiData = new ConcurrentReferenceHashMap<>(20);

    @Autowired
    private TagService tagService;

    @Autowired
    private ImageLinkService imageLinkService;

    @Autowired
    private CacheService cacheService;

    @Autowired
    private TvShowService tvShowService;

    @Override
    public ProgramEntry preProcessItem(ZdfBroadcast broadcast, LocalDate affectedDay) throws ProgramEntryParserException {
        ProgramEntry zdfProgramEntry = new ProgramEntry();

        zdfProgramEntry.setTechnicalId(broadcast.getPosId());
        apiData.put(zdfProgramEntry.getTechnicalId(), broadcast);

        return zdfProgramEntry;
    }

    private void applyStartAndEndDate(ZdfBroadcast broadcast, ProgramEntry programEntry) throws ProgramEntryParserException {
        LocalDateTime startDate;
        LocalDateTime endDate;
        if (broadcast.getEffectiveAirtimeBegin() != null && broadcast.getEffectiveAirtimeEnd() != null) {
            startDate = ZonedDateTime.parse(broadcast.getEffectiveAirtimeBegin()).toLocalDateTime();
            endDate = ZonedDateTime.parse(broadcast.getEffectiveAirtimeEnd()).toLocalDateTime();
        } else {
            startDate = ZonedDateTime.parse(broadcast.getAirtimeBegin()).toLocalDateTime();
            endDate = ZonedDateTime.parse(broadcast.getAirtimeEnd()).toLocalDateTime();
        }

        // detect and set start date + time
        programEntry.setStartDateTime(startDate);
        programEntry.setEndDateTime(endDate);

        // detect and set url of show
        programEntry.setUrl(broadcast.getSelf());

        if (programEntry.getStartDateTime() == null) {
            LOG.fine(programEntry.toString());
            throw new ProgramEntryParserException("start date is null!" + programEntry);
        }
        if (programEntry.getEndDateTime() == null) {
            LOG.fine(programEntry.toString());
            throw new ProgramEntryParserException("end date is null!" + programEntry);
        }
    }

    private String getApiKey() {
        String apiKey = cacheService.getZdfApiKey();
        if (apiKey == null) {
            throw new IllegalStateException("no zdf api key retrievable");
        }
        return apiKey;
    }

    /**
     * this method enriches and finalizes/updates a program entry record
     *
     * @param programEntry
     * @return
     * @throws ProgramEntryParserException
     */
    @Override
    public void postProcessItem(ProgramEntry programEntry) throws ProgramEntryParserException {
        ZdfBroadcast broadcast = apiData.get(programEntry.getTechnicalId());
        if (broadcast == null) {
            throw new ProgramEntryParserException("invalid technical id " + programEntry.getTechnicalId());
        }

        String title = broadcast.getTitle();
        if (broadcast.getSubtitle() != null && !broadcast.getSubtitle().isEmpty()) {
            title += " - " + broadcast.getSubtitle();
        }
        programEntry.setTitle(title);

        applyStartAndEndDate(broadcast, programEntry);

        programEntry.setDescription(broadcast.getText());

        ZdfImage zdfImage = broadcast.getZdfImage();
        if (zdfImage != null ) {

            Layouts zdfLayouts = zdfImage.getLayouts();
            // add image links as list
            List<String> imageUrls = zdfLayouts.getAllImageLinks();
            List<ImageLink> imageLinkList = new ArrayList<>();
            imageUrls.forEach(ele -> {
                LOG.finest(String.format("Storing image url '%s' for program entry", ele));
                imageLinkList.add(imageLinkService.getOrCreateImageLink(ele));
            });
            programEntry.setImageLinks(imageLinkList);
        }
        if (broadcast.getProgrammeItem() != null && !broadcast.getProgrammeItem().isEmpty()) {
            ZdfTagParser tagParser = new ZdfTagParser();

            List<String> tagList = tagParser.getTags(broadcast.getProgrammeItem(), getApiKey());
            List<Tag> tagRecordList = new ArrayList<>();
            tagList.forEach(ele -> {
                LOG.finest(String.format("Store tag '%s'", ele));
                tagRecordList.add(tagService.getOrCreateTag(ele));
            });
            programEntry.setTags(tagRecordList);
            programEntry.setHomePage(ZDFScraper.ZDF_API_HOST + broadcast.getProgrammeItem());
        } else {
            programEntry.setTags(null);
            programEntry.setHomePage(null);
        }
    }

    @Override
    public void linkItem(ProgramEntry programEntry) {
        ZdfBroadcast broadcast = apiData.get(programEntry.getTechnicalId());
        if (broadcast == null) {
            throw new IllegalStateException("invalid technical id " + programEntry.getTechnicalId());
        }
        tvShowService.linkProgramEntryWithZdfBrandId(broadcast.getPrimaryBrandId(), programEntry);
    }

    /**
     * detect elements in html
     *
     * @return
     * @throws ParserException
     */
    @Override
    public Iterable<ZdfBroadcast> getElements(Channel channel, LocalDate day) throws ParserException {
        Iterable<ZdfBroadcast> broadcasts = getRootElement(channel, day);
        if (broadcasts == null) {
            throw new ParserException("No elements retrievable");
        }
        return broadcasts;
    }

    @Override
    public void cleanup() {
        tagService.clear();
    }

    @Override
    public void finishEntry(ProgramEntry programEntry) {
        // free resources in hashmap
        this.apiData.remove(programEntry.getTechnicalId());
    }

    @Override
    public void preProcessProgramList(List<ProgramEntry> linkedProgramList) {
        // do nothing here
    }

    public Iterable<ZdfBroadcast> getRootElement(Channel channel, LocalDate day) {
        final ZDFScraper zdfScraper = new ZDFScraper(channel);
        String queryUrl = zdfScraper.getQueryUrl(day, 1);

        return ZdfApiFetcher.getProgram(queryUrl);
    }

    /**
     * helper inline class around generating ard program urls
     */
    public static class ZDFScraper {

        private static final int LIMIT = 100;
        public static final String DATETIME_SCHEME = "yyyy-MM-dd'T'HH:mm:ssXXX";
        public static final String ZDF_API_HOST = "https://api.zdf.de";
        private Channel channel;
        private final DateTimeFormatter formatter;

        public ZDFScraper() {
            formatter = DateTimeFormatter.ofPattern(DATETIME_SCHEME);
        }

        public ZDFScraper(Channel channel) {
            formatter = DateTimeFormatter.ofPattern(DATETIME_SCHEME);
            this.channel = channel;
        }

        // params: ?datum=&hour=&channel=
        private String ZDF_PROGRAM_URL = ZDF_API_HOST + "/cmdm/epg/broadcasts?limit=" + LIMIT + "&order=asc";

        public Channel getChannel() {
            return channel;
        }

        /**
         * building a url to get channel information
         *
         * @param day
         * @return
         */
        public String getQueryUrl(LocalDate day, int page) {
            if (day == null) {
                throw new IllegalArgumentException("null day object given!");
            }

            LocalDateTime fromDate = LocalDateTime.of(day, LocalTime.of(0,0));
            LocalDateTime toDate = LocalDateTime.of(day.plus(1, ChronoUnit.DAYS), LocalTime.of(0,0));

            StringBuilder argumentList = new StringBuilder();
            if (channel != null) {
                argumentList.append("&tvServices=").append(channel.getTechnicalId());
            }
            argumentList.append("&from=").append(getIsoDate(fromDate));
            argumentList.append("&to=").append(getIsoDate(toDate));

            if (page != 0) {
                argumentList.append("&page=" + page);
            }

            // build url
            return ZDF_PROGRAM_URL + argumentList.toString();
        }

        private String getIsoDate(LocalDateTime day) {
            try {
                return URLEncoder.encode(formatter.format(ZonedDateTime.of(day, ZoneId.of("Europe/Berlin"))), StandardCharsets.UTF_8.toString());
            } catch (UnsupportedEncodingException e) {
                LOG.throwing(ProgramEntryParser.ZDFScraper.class.getName(), "getIsoDate", e);
            }
            return null;
        }

        /**
         * method to fetch zdf api key from zdf home page
         * @return
         */
        public String retrieveApiKey() {
            LOG.info("Start retrieving zdf api key");
            Document jsoupDoc = Fetcher.fetchDocument(ZDF_HOST + "/live-tv", "body");
            Elements elements = jsoupDoc.select("script");
            String apiKey = null;
            for (Element ele : elements) {
                final String html = ele.html();
                if(html.contains("IMPORTANT CONFIGURATION!")) {
                    Pattern pattern = Pattern.compile("apiToken: '(.*?)'");
                    Matcher matcher = pattern.matcher(html);
                    if (matcher.find()) {
                        apiKey = matcher.group()
                                .replace("apiToken:", "")
                                .replace("'", "").trim();
                    }
                    break;
                }
            }
            return apiKey;
        }
    }

    public static class ZdfTagParser {
        public List<String> getTags(String zdfPosUrl, String apiKey) {
            if (zdfPosUrl == null || zdfPosUrl.isEmpty()) {
                throw new IllegalArgumentException("invalid null or empty url given");
            }
            zdfPosUrl = ZDFScraper.ZDF_API_HOST + zdfPosUrl;
            List<String> tagList = new ArrayList<>();
            TvShowModel tvShowModel = ZdfApiFetcher.getTvShow(zdfPosUrl);

            if (tvShowModel.getCategory() != null && !tvShowModel.getCategory().isEmpty()) {
                tagList.add(tvShowModel.getCategory());
            }
            if (tvShowModel.getCategory() == null) {
                LOG.finest("zdf show '" + tvShowModel.getTitle() + "' has no category in api response. url: " + zdfPosUrl);
            }
            if (tvShowModel.getGenre() != null && !tvShowModel.getGenre().isEmpty()
                && tvShowModel.getCategory() != null && !tvShowModel.getCategory().equals(tvShowModel.getGenre())) {
                tagList.add(tvShowModel.getGenre());
            }
            return tagList;
        }
    }
}
