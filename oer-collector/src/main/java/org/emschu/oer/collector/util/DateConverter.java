package org.emschu.oer.collector.util;

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

import javax.validation.constraints.NotNull;
import java.time.DateTimeException;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.logging.Logger;

/**
 * class to convert given dates/times in orf parsing context
 */
public class DateConverter {

    private static final Logger LOG = Logger.getLogger(DateConverter.class.getName());

    /**
     * method to parse a time string and combine it with the correct day to a {@link LocalDateTime} object
     *
     * @param isNight
     * @param day
     * @param timeString
     * @return
     */
    public static LocalDateTime generateDateForEntry(boolean isNight,
                                                     @NotNull LocalDate day,
                                                     @NotNull String timeString) {
        final int middle = timeString.indexOf(':');
        if (middle == -1) {
            LOG.warning("invalid time string given: " + timeString);
            return null;
        }

        final int hour;
        final int minutes;
        try {
            hour = Integer.valueOf(timeString.substring(Math.max(middle - 2, 0), middle));
            minutes = Integer.valueOf(timeString.substring(middle + 1, middle + 3));
            // check if second number = minutes are not split
            if (timeString.length() >= middle + 4) {
                if (Integer.parseInt(String.valueOf(timeString.charAt(middle + 4))) >= 0) {
                    return null;
                }
            }
        } catch (NumberFormatException nfe) {
            LOG.warning("invalid time string given: " + timeString);
            return null;
        }
        if (hour > 60 || minutes > 60) {
            // invalid values!
            return null;
        }

        final LocalDateTime ret;
        try {
            if (isNight && hour < 8) {
                // add one day
                ret = LocalDateTime.of(day.plus(1, ChronoUnit.DAYS), LocalTime.of(hour, minutes));
            } else {
                ret = LocalDateTime.of(day, LocalTime.of(hour, minutes));
            }
        } catch (DateTimeException dte) {
            LOG.warning("invalid date time values in: " + timeString);
            return null;
        }
        return ret;
    }

    /**
     * made for strings like:
     *  Dienstag, 09.04.2019
     *
     * TODO test
     * @param dateString
     * @return
     */
    public static LocalDate getDateInString(String dateString) {
        if (dateString == null || dateString.isEmpty()) {
            return null;
        }
        final int firstPointIndex = dateString.indexOf('.');
        if (firstPointIndex == -1) {
            // not a valid string
            return null;
        }
        final int secondPointIndex = dateString.indexOf('.', firstPointIndex + 1);
        if (secondPointIndex == -1) {
            // we need two points '.' in the string to work
            return null;
        }
        try {
            if (dateString.length() < secondPointIndex + 5) {
                // not plausible
                return null;
            }
            final int day = Integer.parseInt(dateString.substring(firstPointIndex - 2, firstPointIndex));
            final int month = Integer.parseInt(dateString.substring(firstPointIndex + 1, secondPointIndex));
            final int year = Integer.parseInt(dateString.substring(secondPointIndex + 1, secondPointIndex + 5));
            return LocalDate.of(year, month, day);
        } catch (NumberFormatException | DateTimeException nfe) {
            LOG.warning("invalid dateString in string: " + dateString);
            return null;
        }
    }
}
