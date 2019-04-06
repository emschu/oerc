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

import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.collector.service.ProgramService;
import org.emschu.oer.collector.service.TagService;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.Tag;
import org.emschu.oer.core.util.Hasher;
import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.logging.Logger;

@Component("orfProgramEntryParser")
public class ProgramEntryParser implements ProgramEntryParserInterface<Element> {

    private static final Logger LOG = Logger.getLogger(ProgramEntryParser.class.getName());

    @Autowired
    private TagService tagService;

    @Autowired
    private ProgramService programService;

    @Override
    public ProgramEntry preProcessItem(Element element, LocalDate affectedDay) throws ProgramEntryParserException {
        final Elements title = element.select("h2");
        final Elements subTitle = element.select("div.teaser h3");
        final Elements detailLink = element.select("p.detaillink");
        final Elements startTime = element.select("div.starttime h3");
        final Elements description = element.select("div.storytext");
        final boolean isNightEntry = element.className().contains("tsnight");

        final String entryTitle;
        if (!subTitle.isEmpty()) {
            entryTitle = title.text() + " - " + subTitle.text();
        } else {
            entryTitle = title.text();
        }

        final String url = detailLink.select("a").attr("href");
        if (entryTitle.isEmpty()) {
            LOG.warning("no title found. Skip entry.");
            // leave
            return null;
        }
        if (startTime.isEmpty()) {
            LOG.warning("skip entry. no start time for title " + entryTitle);
            return null;
        }

        ProgramEntry programEntry = new ProgramEntry();
        if (!url.isEmpty()) {
            programEntry.setUrl(url);
        }

        final LocalDateTime startDateTime = ORFDateConverter.generateDateForEntry(isNightEntry, affectedDay, startTime.text());
        if (startDateTime == null) {
            LOG.warning("something went wrong parsing: " + startTime.text());
            return null;
        }
        programEntry.setStartDateTime(startDateTime);
        programEntry.setTitle(entryTitle);
        if (!description.isEmpty()) {
            programEntry.setDescription(description.text());
        }

        final String uniqueHash = startDateTime.format(DateTimeFormatter.ISO_LOCAL_DATE_TIME) + ":"
                                    + entryTitle + ":" + url;
        LOG.fine("unique hash for program entry: " + uniqueHash);
        programEntry.setTechnicalId(Hasher.getHash(uniqueHash));

        return programEntry;
    }

    @Override
    public void postProcessItem(ProgramEntry programEntry) throws ProgramEntryParserException {
        final String url = programEntry.getUrl();
        if (url == null) {
            LOG.fine("No further information for program entry");
            return;
        }
        if (url.contains("okidoki.orf.at")) {
            // quiet skipping of these pages - there is too less information
            return;
        }
        LOG.fine("calling url: " + url);

        final Document page = getPage(url);
        final String fullDescription = page.select(".paragraph").text();

        if (fullDescription == null || fullDescription.isEmpty()) {
            LOG.fine("Empty description for url: " + url);
        } else {
            programEntry.setDescription(fullDescription);
        }

        final String homePage = page.select("div.status p.network a").attr("href");
        if (homePage != null && !homePage.isEmpty()) {
            programEntry.setHomePage(homePage);
        }

        // store genre as tag
        final String genre = page.select("div.starttime p.genre").text();
        if (genre != null && !genre.isEmpty()) {
            ArrayList<Tag> tagRecordList = new ArrayList<>();
            LOG.fine(String.format("Store tag '%s'", genre));
            tagRecordList.add(tagService.getOrCreateTag(genre.trim()));
            programEntry.setTags(tagRecordList);
        } else {
            LOG.fine("No tag found for url: " + url);
        }
    }

    @Override
    public void linkItem(ProgramEntry programEntry) { }

    @Override
    public Iterable<Element> getElements(Channel channel, LocalDate day) throws ParserException {
        final Elements rootElement = getRootElement(channel, day);
        if (rootElement == null) {
            LOG.warning("No orf root element for day: " + day.format(DateTimeFormatter.ISO_LOCAL_DATE));
            return new Elements();
        }

        Document jsoupDoc = Jsoup.parse(rootElement.outerHtml());
        Elements foundElements = jsoupDoc.select("div.clearer");
        if (foundElements == null) {
            throw new ParserException("No elements retrievable");
        }
        return foundElements;
    }

    private String ensureTwoDigits(int digit) {
        if (digit < 10) {
            return "0" + digit;
        }
        return String.valueOf(digit);
    }

    private Elements getRootElement(Channel channel, LocalDate day) {
        Elements mainElementDiv = null;

        // future: 23 days ?
        // past: 15 days ?
        long daysBetween = ChronoUnit.DAYS.between(day, LocalDate.now());
        // negative values represent future
        if (daysBetween > 14 || daysBetween < -22) {
            LOG.warning("Invalid date range for orf data request.");
            return null;
        }

        String queryUrl = "https://tv.orf.at/program/" + channel.getTechnicalId() + "/" + day.getYear() +
                ensureTwoDigits(day.getMonthValue())
                + ensureTwoDigits(day.getDayOfMonth());

        Document jsoupDoc = getPage(queryUrl);
        mainElementDiv = jsoupDoc.body().select("div.main");
        if (mainElementDiv == null || mainElementDiv.size() != 1) {
            // we expect exactly one element with this class
            throw new IllegalStateException("no event list tag in page html found");
        }
        return mainElementDiv;
    }

    private Document getPage(String queryUrl) {
        LOG.fine("Query url: " + queryUrl);

        // build orf headers
        Map<String, String> headers = new HashMap<>();
        headers.put("Origin", "tv.orf.at");
        headers.put("Host", "tv.orf.at");

        return Fetcher.fetchDocument(queryUrl, "body", headers);
    }

    @Override
    public void cleanup() { }

    @Override
    public void finishEntry(ProgramEntry programEntry) { }

    @Override
    public void preProcessProgramList(List<ProgramEntry> linkedProgramList) {
        tryToDetectEndDates(linkedProgramList);

        if (linkedProgramList.isEmpty()) {
            return;
        }

        // we handle an orf particularity here:
        // there is no end date time in orf web pages. we can set the end date time to (n+1).startDateTime,
        // as the method above is implemented.
        // This works for all items, except the last
        final ProgramEntry lastProgramEntry = linkedProgramList.get(linkedProgramList.size() - 1);
        if (lastProgramEntry.getEndDateTime() == null) {
            final Optional<ProgramEntry> nearestEntry = programService.nearestProgramEntryInFuture(lastProgramEntry.getStartDateTime(), lastProgramEntry.getChannel());
            if (nearestEntry.isPresent()) {
                LOG.fine("nearest entry is: " + nearestEntry.get());
                long plausibilityCheck = ChronoUnit.MINUTES.between(lastProgramEntry.getStartDateTime(), nearestEntry.get().getStartDateTime());
                if (plausibilityCheck < 250) {
                    lastProgramEntry.setEndDateTime(nearestEntry.get().getStartDateTime());
                } else {
                    LOG.warning("nearest entry is not plausible");
                }
            } else {
                LOG.warning("No nearest entry for program entry: " + lastProgramEntry);
            }
        }
    }
}
