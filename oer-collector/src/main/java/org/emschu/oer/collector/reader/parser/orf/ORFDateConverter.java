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
public class ORFDateConverter {

    private static final Logger LOG = Logger.getLogger(ORFDateConverter.class.getName());

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
            hour = Integer.valueOf(timeString.substring(0, middle));
            minutes = Integer.valueOf(timeString.substring(middle + 1, middle + 3));
        } catch (NumberFormatException nfe) {
            LOG.warning("invalid time string given: " + timeString);
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
}
