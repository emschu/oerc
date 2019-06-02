package org.emschu.oer.oerserver;

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

import org.emschu.oer.collector.util.DateConverter;
import org.junit.Assert;
import org.junit.Test;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;

public class DateConverterTest {

    @Test
    public void testStartTimeParsing() {
        final LocalDateTime localDateTime1 = DateConverter.generateDateForEntry(false, LocalDate.now(), "7:00");
        final LocalDateTime localDateTime2 = DateConverter.generateDateForEntry(true, LocalDate.now(), "21:50");
        final LocalDateTime localDateTime3 = DateConverter.generateDateForEntry(false, LocalDate.now(), "21:50");
        final LocalDateTime localDateTime4 = DateConverter.generateDateForEntry(true, LocalDate.now(), "5:00");

        Assert.assertEquals(LocalDateTime.of(LocalDate.now(), LocalTime.of(7,0)), localDateTime1);
        Assert.assertEquals(LocalDateTime.of(LocalDate.now(), LocalTime.of(21,50)), localDateTime2);
        Assert.assertEquals(LocalDateTime.of(LocalDate.now(), LocalTime.of(21,50)), localDateTime3);
        Assert.assertEquals(LocalDateTime.of(LocalDate.now().plus(1, ChronoUnit.DAYS),
                LocalTime.of(5,0)), localDateTime4);
    }

    @Test
    public void testWrongResult() {
        Assert.assertNull(DateConverter.generateDateForEntry(true, LocalDate.now(), ""));
        Assert.assertNull(DateConverter.generateDateForEntry(true, LocalDate.now(), "test:test"));
        Assert.assertNull(DateConverter.generateDateForEntry(true, LocalDate.now(), "500:5000"));
        Assert.assertNull(DateConverter.generateDateForEntry(true, null, ""));
    }

    @Test
    public void testDateOfStringExtraction() {
        final LocalDate dateInString = DateConverter.getDateInString("Dienstag, 02.03.2019");
        Assert.assertNotNull(dateInString);
        Assert.assertEquals(2, dateInString.getDayOfMonth());
        Assert.assertEquals(3, dateInString.getMonthValue());
        Assert.assertEquals(2019, dateInString.getYear());
    }

    @Test
    public void testDateOfStringExtractionFails() {
        Assert.assertNull(DateConverter.getDateInString(null));
        Assert.assertNull(DateConverter.getDateInString(""));
        Assert.assertNull(DateConverter.getDateInString("02032019"));
        Assert.assertNull(DateConverter.getDateInString("02.011234"));
        Assert.assertNull(DateConverter.getDateInString("Montag Dienstag"));
        Assert.assertNull(DateConverter.getDateInString("Dolby 5.1"));
    }
}
