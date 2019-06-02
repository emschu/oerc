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

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.collector.service.ImageLinkService;
import org.emschu.oer.collector.service.TagService;
import org.emschu.oer.collector.util.EidExtractor;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ImageLink;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.Tag;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.util.Pair;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.logging.Logger;

/**
 * a concrete implementation of {@link ProgramEntryParserInterface} to fetch and enrich ard program data
 */
@Component(value = "ardProgramEntryParser")
public class ProgramEntryParser implements ProgramEntryParserInterface<Element> {

    public static final String ARD_HOST = "https://programm.ard.de";
    private static final Logger LOG = Logger.getLogger(ProgramEntryParser.class.getName());

    @Autowired
    private TagService tagService;

    @Autowired
    private ImageLinkService imageLinkService;

    @Override
    public ProgramEntry preProcessItem(Element html, LocalDate affectedDay, Channel channel) throws ProgramEntryParserException {
        ProgramEntry ardProgramEntry = new ProgramEntry();
        Document jsoupDoc = Jsoup.parse(html.outerHtml());

        // handle title
        String title = jsoupDoc.select("span.title").text();
        ardProgramEntry.setTitle(title);

        // detect and set start date + time
        String startDateTimeStr = jsoupDoc.select("span.date").text();
        int hour = Integer.parseInt(startDateTimeStr.substring(0, startDateTimeStr.indexOf(':')));
        int minute = Integer.parseInt(startDateTimeStr.substring(startDateTimeStr.indexOf(':') + 1, startDateTimeStr.length()));

        ardProgramEntry.setStartDateTime(LocalDateTime.of(affectedDay, LocalTime.of(hour, minute)));

        // detect and set technical id
        String technicalIdRaw = jsoupDoc.getElementsByAttributeValueContaining("class", "eid").attr("class");
        ardProgramEntry.setTechnicalId(EidExtractor.extractEid(technicalIdRaw));

        // detect and set url of show
        String urlOfShow = jsoupDoc.getElementsByAttribute("href").attr("href");
        ardProgramEntry.setUrl(urlOfShow);

        return ardProgramEntry;
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
        Elements relevantDiv = fetchProgramPage(programEntry.getUrl());
        if (relevantDiv == null) {
            throw new ProgramEntryParserException("invalid root element of program entry detail");
        }
        // start and end field is updated
        applyICalData(relevantDiv, programEntry);

        if (programEntry.getStartDateTime() == null) {
            LOG.fine(programEntry.toString());
            throw new ProgramEntryParserException("start date is null!" + programEntry);
        }
        if (programEntry.getEndDateTime() == null) {
            LOG.fine(programEntry.toString());
            throw new ProgramEntryParserException("end date is null!" + programEntry);
        }
        applyHtmlData(relevantDiv, programEntry);
    }

    @Override
    public void linkItem(ProgramEntry programEntry) {

    }

    private Elements fetchProgramPage(String url) {
        String programEntryPage = ARD_HOST + url;
        LOG.fine(String.format("parsing page: '%s'", programEntryPage));
        Document jsoupDoc = Fetcher.fetchDocument(programEntryPage, "body");
        return jsoupDoc.body().select(".program-con");
    }

    public void applyHtmlData(Elements relevantDiv, ProgramEntry programEntry) {
        // detect title
        String titleInformationHTML = relevantDiv.select("span.title").html();
        titleInformationHTML = titleInformationHTML.replace("<span class=\"subtitle\">", "-")
                .replace("</span>", "").trim();
        if (titleInformationHTML.contains("|")) {
            String title = titleInformationHTML.substring(0, titleInformationHTML.indexOf('|')).trim();
            programEntry.setTitle(title);
        }
        // add description text
        String descriptionText = parseDescription(relevantDiv, programEntry.getTechnicalId(), programEntry.getUrl());
        programEntry.setDescription(descriptionText);

        // add image links as list
        Elements images = relevantDiv.select(".media-con img");
        List<String> imageUrls = new ArrayList<>();
        images.forEach(element -> {
            imageUrls.add(element.attr("src"));
        });
        ArrayList<ImageLink> imageLinkList = new ArrayList<>();
        imageUrls.forEach(ele -> {
            LOG.finest(String.format("Storing image url '%s' for program entry", ele));
            imageLinkList.add(imageLinkService.getOrCreateImageLink(ele));
        });
        programEntry.setImageLinks(imageLinkList);

        // get genre data
        ArdTagParser ardTagParser = new ArdTagParser();
        List<String> tagList = ardTagParser.getTags(programEntry.getTechnicalId());
        ArrayList<Tag> tagRecordList = new ArrayList<>();
        tagList.forEach(ele -> {
            LOG.finest(String.format("Store tag '%s'", ele));
            tagRecordList.add(tagService.getOrCreateTag(ele.trim()));
        });
        programEntry.setTags(tagRecordList);

        // detect "real" home page of this program entry - if any
        String hp = getHomePage(relevantDiv);
        programEntry.setHomePage(hp);
    }

    /**
     * detect special homepage of a tv show
     *
     * @param relevantDiv
     * @return
     */
    private String getHomePage(Elements relevantDiv) {
        Elements links = relevantDiv.select(".bcData").select("a");
        for (Element ele : links) {
            if (ele.html().contains("Sendungsseite im Internet")) {
                return ele.attr("href");
            }
        }
        return null;
    }

    /**
     * parse description of a program entry
     *
     * @param relevantDiv
     * @param technicalId
     * @return
     */
    private String parseDescription(Elements relevantDiv, String technicalId, String url) {
        Elements descElements = relevantDiv.select(String.format("#mehr-%s .eventText", technicalId));
        String descriptionText = descElements.text();
        if (descriptionText.contains("Keine weiteren Informationen") || descriptionText.isEmpty()) {
            descriptionText = null;
        }
        if (descriptionText == null) {
            // use alternative.. - if existent
            if (relevantDiv.select(".detail-top").isEmpty()) {
                LOG.warning("Fallback mechanism of description text fetching. Re-fetch data from '" + url + "'");
                // fallback mechanism
                relevantDiv = fetchProgramPage(url);
            }
            Elements eventText = relevantDiv.select("div.detail-top").select("div.eventText");
            if (eventText.isEmpty()) {
                descriptionText = null;
            } else {
                descriptionText = eventText.get(0).text();
            }
        }

        if (descriptionText == null || descriptionText.isEmpty()) {
            return null;
        }
        return descriptionText;
    }

    /**
     * central method to update program entry start and end dates with the
     * assumed the most correct information about start and end of a single tv show.
     *
     * @param html
     * @param programEntry
     */
    public void applyICalData(Elements html, ProgramEntry programEntry) {
        if (programEntry == null) {
            throw new IllegalArgumentException("null entry given");
        }
        ArdICalParser ardIcalParser = new ArdICalParser();
        String icalLink = ardIcalParser.parseIcalLink(html);
        if (icalLink == null || icalLink.isEmpty()) {
            LOG.warning("invalid ical link given");
            return;
        }

        String icalContent = ardIcalParser.getFileContent(icalLink);
        if (icalContent == null || icalContent.isEmpty()) {
            LOG.warning("invalid ical content given");
            return;
        }
        Pair<LocalDateTime, LocalDateTime> dateValues = ardIcalParser.getRelevantDates(icalContent);
        if (dateValues != null) {
            programEntry.setStartDateTime(dateValues.getFirst());
            programEntry.setEndDateTime(dateValues.getSecond());
        } else {
            LOG.warning("no start and end information retrievable from ical content");
            LOG.finest(programEntry.toString());
        }
    }

    /**
     * detect elements in html
     *
     * @return
     * @throws ParserException
     */
    @Override
    public Iterable<Element> getElements(Channel channel, LocalDate day) throws ParserException {
        Document jsoupDoc = Jsoup.parse(getRootElement(channel, day).outerHtml());

        Elements foundElements = jsoupDoc.select("li[class^=eid]");
        if (foundElements == null) {
            throw new ParserException("No elements retrievable");
        }
        return foundElements;
    }

    @Override
    public void cleanup() {
        tagService.clear();
    }

    @Override
    public void finishEntry(ProgramEntry programEntry) {
        // do nothing here
    }

    @Override
    public void preProcessProgramList(List<ProgramEntry> linkedProgramList) {
        tryToDetectEndDates(linkedProgramList);
    }

    public Elements getRootElement(Channel channel, LocalDate day) {
        Elements eventListTag = null;
        String queryUrl = new ARDScraper(channel).getQueryUrl(day);
        Document jsoupDoc = Fetcher.fetchDocument(queryUrl, ".event-list");
        eventListTag = jsoupDoc.body().getElementsByClass("event-list");
        if (eventListTag == null || eventListTag.size() != 1) {
            // we expect exactly one element with this class
            throw new IllegalStateException("no event list tag in page html found");
        }
        return eventListTag;
    }

    /**
     * helper inline class around generating ard program urls
     */
    public static class ARDScraper {

        private Channel channel;

        public ARDScraper(Channel channel) {
            this.channel = channel;
        }

        // params: ?datum=&hour=&channel=
        private String ARD_PROGRAM_URL = ARD_HOST + "/TV/Programm/Sender";

        public Channel getChannel() {
            return channel;
        }

        /**
         * building a url to get channel information
         *
         * @param day
         * @return
         */
        public String getQueryUrl(LocalDate day) {
            if (day == null) {
                throw new IllegalArgumentException("null day object given!");
            }
            // build url
            return ARD_PROGRAM_URL + "?datum=" + day.getDayOfMonth() + "." + day.getMonthValue() + "." + day.getYear() +
                    "&hour=0&sender=" + channel.getTechnicalId();
        }
    }

    /**
     * parses start and end date
     * TODO Test + extract
     */
    public static class ArdTagParser {
        private static final String TAG_URL = "/TV/Programm/Load/Similar35?eid=";

        public List<String> getTags(String technicalId) {
            List<String> tagList = new ArrayList<>();
            Document jsoupDoc = Fetcher.fetchDocument(ARD_HOST + TAG_URL + technicalId, "body");
            jsoupDoc.select("form[id^=bookmark-checks] .row span[class*=similar-events-bookmark]")
                    .forEach(element -> tagList.add(element.text()));
            return tagList;
        }

    }
}
