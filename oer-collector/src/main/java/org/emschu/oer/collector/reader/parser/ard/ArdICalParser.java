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

import net.fortuna.ical4j.data.CalendarBuilder;
import net.fortuna.ical4j.model.*;
import net.fortuna.ical4j.model.component.CalendarComponent;
import net.fortuna.ical4j.model.property.CalScale;
import net.fortuna.ical4j.model.property.DtEnd;
import net.fortuna.ical4j.model.property.DtStart;
import net.fortuna.ical4j.model.property.TzId;
import org.jsoup.select.Elements;
import org.springframework.data.util.Pair;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.StringReader;
import java.net.MalformedURLException;
import java.net.URL;
import java.time.LocalDateTime;
import java.util.logging.Logger;

/**
 * wrapper around ical4j library to extract start and end date of an ard program entry
 * has no real internal state
 */
public class ArdICalParser {

    private static final Logger LOG = Logger.getLogger(ArdICalParser.class.getName());
    private final CalendarBuilder builder = new CalendarBuilder();
    private final TimeZoneRegistry registry;

    public ArdICalParser() {
        registry = builder.getRegistry();
    }

    /**
     * parses an ical link - if found or returns null
     *
     * @param html jsoup object
     * @return or null on failure
     */
    public String parseIcalLink(Elements html) {
        String icalLink = html.select("a[href^=/ICalendar]").attr("href");
        if (icalLink == null || icalLink.isEmpty()) {
            LOG.warning("No ical link in html found");
            return null;
        }
        return icalLink;
    }

    /**
     * parses a link to string content
     *
     * @return or null on failure
     */
    public String getFileContent(String icalLink) {
        StringBuilder icalFile = new StringBuilder();
        URL icalUrl = null;
        try {
            icalUrl = new URL("https", ProgramEntryParser.ARD_HOST.replace("https://", ""), icalLink);
        } catch (MalformedURLException e) {
            LOG.throwing(ArdICalParser.class.getName(), "getFileContent", e);
            return null;
        }
        // reading line by line
        try (BufferedReader bis = new BufferedReader(new InputStreamReader(icalUrl.openStream()))) {
            String str = null;
            while ((str = bis.readLine()) != null) {
                icalFile.append(str).append(System.lineSeparator());
            }
        } catch (IOException e) {
            LOG.throwing(ArdICalParser.class.getName(), "getFileContent", e);
            return null;
        }
        return icalFile.toString();
    }

    /**
     * method to find start and end date out of an ical file
     *
     * @param icalContent
     * @return
     */
    public Pair<LocalDateTime, LocalDateTime> getRelevantDates(String icalContent) {
        if (icalContent == null) {
            throw new IllegalArgumentException("null ical content given");
        }
        // fix ics format given by ard
        icalContent = applyFixesToICalFormat(icalContent);
        try {
            final StringReader in = new StringReader(icalContent);
            Calendar calendar = builder.build(in);
            in.close();
            calendar.getProperties().add(CalScale.GREGORIAN);
            CalendarComponent vtimezone = calendar.getComponent("VTIMEZONE");
            if (vtimezone == null) {
                LOG.warning("no VTIMEZONE element found in ical content");
                return null;
            }
            TzId tzid = vtimezone.getProperty("TZID");
            if (tzid != null) {
                TimeZone tz = registry.getTimeZone(tzid.getValue());
                for (Component calComponent : calendar.getComponents("VEVENT")) {
                    DtStart startDateFromIcal = calComponent.getProperty(Property.DTSTART);
                    DtEnd endDateFromIcal = calComponent.getProperty(Property.DTEND);
                    LocalDateTime calendarStartDate = LocalDateTime.ofInstant(startDateFromIcal.getDate().toInstant(), tz.toZoneId());
                    LocalDateTime calendarEndDate = LocalDateTime.ofInstant(endDateFromIcal.getDate().toInstant(), tz.toZoneId());

                    return Pair.of(calendarStartDate, calendarEndDate);
                }
                LOG.fine("Successfully parsed ard ics content");
            } else {
                LOG.warning("No proper time zone given in ical calendar");
            }
        } catch (IOException | net.fortuna.ical4j.data.ParserException e) {
            LOG.finest("content: " + icalContent);
            LOG.throwing(ProgramEntryParser.class.getName(), "filterDetailInformation", e);
            LOG.warning("problem fetching start and end date: " + e.getMessage());
        }
        return null;
    }

    /**
     * the generated ical of ard does not fit to ical specs..
     *
     * @param icalContent
     * @return
     */
    public String applyFixesToICalFormat(String icalContent) {
        final String shittyTzString = "TZID=Europe/Berlin:";
        icalContent = icalContent.replace("DTSTART;" + shittyTzString, "DTSTART:");
        icalContent = icalContent.replace("DTEND;" + shittyTzString, "DTEND:");
        icalContent = icalContent.replace("DTSTAMP;" + shittyTzString, "DTSTAMP:");
        return icalContent;
    }
}
