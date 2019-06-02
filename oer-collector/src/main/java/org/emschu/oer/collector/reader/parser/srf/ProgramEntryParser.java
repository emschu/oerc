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

import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserException;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.collector.service.TagService;
import org.emschu.oer.collector.util.DateConverter;
import org.emschu.oer.collector.util.StringFormat;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.Tag;
import org.emschu.oer.core.util.Hasher;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Logger;

@Component("srfProgramEntryParser")
public class ProgramEntryParser implements ProgramEntryParserInterface<Element> {

    private static final Logger LOG = Logger.getLogger(ProgramEntryParser.class.getName());

    @Autowired
    private TagService tagService;

    @Override
    public ProgramEntry preProcessItem(Element element, LocalDate affectedDay, Channel channel) throws ProgramEntryParserException {
        final String startDate = element.select(".channel-show__airtime .channel-show__begin").text();
        final String endDate = element.select(".channel-show__airtime .channel-show__stop").text();
        final String url = element.select(".channel-show__link").attr("href");
        final String title = element.select(".channel-show .channel-show__title").text();
        ProgramEntry programEntry = new ProgramEntry();
        // in pre-processing we cannot decide on which day the given entry is televised. #
        // we do this in post processing (finally)
        if (startDate != null) {
            programEntry.setStartDateTime(DateConverter.generateDateForEntry(false, affectedDay, startDate));
        }
        if (endDate != null) {
            programEntry.setEndDateTime(DateConverter.generateDateForEntry(false, affectedDay, endDate));
        }
        if (programEntry.getStartDateTime() == null) {
            LOG.warning("no start date time for program entry found!");
            return null;
        }
        programEntry.setTitle(title);
        if (url != null) {
            programEntry.setUrl(url.replaceFirst("//", "https://"));
        }

        final String uniqueHash = programEntry.getStartDateTime().format(DateTimeFormatter.ISO_LOCAL_DATE_TIME) + ":"
                + title + ":" + url + ":" + channel.getId();
        programEntry.setTechnicalId(Hasher.getHash(uniqueHash));

        return programEntry;
    }

    @Override
    public void postProcessItem(ProgramEntry programEntry) throws ProgramEntryParserException {
        final String detailUrl = programEntry.getUrl();
        if (detailUrl == null || detailUrl.isEmpty()) {
            LOG.info("No url for data" + programEntry.toString());
            return;
        }
        LOG.fine(String.format("Page detail url: '%s'", detailUrl));

        final Document body = Fetcher.fetchDocument(detailUrl, "body #content");
        final Elements infoboxParagraphs = body.select(".detail--content .infobox .left p");
        if (infoboxParagraphs.size() > 2) {
            // date + time
            final String date = infoboxParagraphs.get(1).text() + " " + infoboxParagraphs.get(0).text();
            LocalDate realDateOfEntry = DateConverter.getDateInString(date);

            final LocalDateTime initialStartDateTime = programEntry.getStartDateTime();
            if (initialStartDateTime.getDayOfMonth() != realDateOfEntry.getDayOfMonth()) {
                // adjust night dates, store duration
                long durationInMinutes = ChronoUnit.MINUTES.between(initialStartDateTime, programEntry.getEndDateTime());

                programEntry.setStartDateTime(LocalDateTime.of(realDateOfEntry,
                        LocalTime.of(initialStartDateTime.getHour(), initialStartDateTime.getMinute(), initialStartDateTime.getSecond())));

                LocalDateTime newEndDateTime = programEntry.getStartDateTime().plus(durationInMinutes, ChronoUnit.MINUTES);
                programEntry.setEndDateTime(newEndDateTime);

                LOG.fine(String.format("Adjust date of program entry from '%s' to '%s'",
                        initialStartDateTime, programEntry.getStartDateTime()));
            }
            if (programEntry.getEndDateTime().isBefore(programEntry.getStartDateTime())) {
                LocalDateTime shiftEndDateTime = programEntry.getEndDateTime().plus(1, ChronoUnit.DAYS);
                programEntry.setEndDateTime(shiftEndDateTime);
                LOG.fine(String.format("Adjust enddate of program entry to '%s'", programEntry.getEndDateTime()));
            }
        }

        final String description = extractDescription(body);
        if (!description.isEmpty()) {
            programEntry.setDescription(description);
        }

        if (infoboxParagraphs.size() > 4) {
            String genre = null;
            if (infoboxParagraphs.get(2).text().contains("Wiederholung")) {
                genre = infoboxParagraphs.get(3).text();
            } else {
                genre = infoboxParagraphs.get(2).text();
            }
            if (genre != null) {
                // store genre
                ArrayList<Tag> tagRecordList = new ArrayList<>();
                LOG.fine(String.format("Store tag '%s'", genre));
                tagRecordList.add(tagService.getOrCreateTag(genre.trim()));
                programEntry.setTags(tagRecordList);
            } else {
                LOG.warning(String.format("No genre for program entry in url '%s'", detailUrl));
            }
        }
    }

    /**
     * helper method to extract description
     *
     * @param body
     * @return
     */
    private String extractDescription(Document body) {
        final Elements leadElements = body.select("p.lead");
        final Elements descElements = body.select("p.description");
        final Elements actorElements = body.select("ul.actors");

        StringBuilder descBuilder = new StringBuilder();
        if (!leadElements.isEmpty()) {
            descBuilder.append(leadElements.text()).append("\n");
        }
        if (!descElements.isEmpty()) {
            descBuilder.append(descElements.text()).append("\n");
        }
        if (!actorElements.isEmpty()) {
            descBuilder.append(actorElements.text()).append("\n");
        }
        return descBuilder.toString();
    }

    @Override
    public void linkItem(ProgramEntry programEntry) {
        // do nothing here for orf
    }

    @Override
    public Iterable<Element> getElements(Channel channel, LocalDate day) throws ParserException {
        // future: 29 days
        // past: 15 days
        long daysBetween = ChronoUnit.DAYS.between(day, LocalDate.now());
        // negative values represent future
        if (daysBetween > 15 || daysBetween < -30) {
            LOG.warning("Invalid date range for srf data request.");
            return null;
        }

        final String url = "https://www.srf.ch/programm/tv/sender/" + channel.getTechnicalId() + "/"
                + StringFormat.ensureTwoDigits(day.getDayOfMonth()) + "-"
                + StringFormat.ensureTwoDigits(day.getMonthValue()) + "-"
                + StringFormat.ensureTwoDigits(day.getYear());

        final Document body = Fetcher.fetchDocument(url, "body");
        if (body == null) {
            LOG.warning(String.format("No content in url '%s'", url));
            return null;
        }

        final Elements select = body.select("#content .channel-show");
        if (select == null) {
            throw new ParserException("no entries found for url: " + url);
        }
        return select;
    }

    @Override
    public void cleanup() {
        // do nothing here for orf
    }

    @Override
    public void finishEntry(ProgramEntry programEntry) {
        // do nothing here for orf
    }

    @Override
    public void preProcessProgramList(List<ProgramEntry> linkedProgramList) {
        // do nothing here
    }
}
