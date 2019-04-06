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

import org.emschu.oer.collector.reader.parser.orf.ORFDateConverter;
import org.junit.Assert;
import org.junit.Test;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;

public class OrfDateConverterTest {

    @Test
    public void testStartTimeParsing() {
        final LocalDateTime localDateTime1 = ORFDateConverter.generateDateForEntry(false, LocalDate.now(), "7:00");
        final LocalDateTime localDateTime2 = ORFDateConverter.generateDateForEntry(true, LocalDate.now(), "21:50");
        final LocalDateTime localDateTime3 = ORFDateConverter.generateDateForEntry(false, LocalDate.now(), "21:50");
        final LocalDateTime localDateTime4 = ORFDateConverter.generateDateForEntry(true, LocalDate.now(), "5:00");

        Assert.assertEquals(LocalDateTime.of(LocalDate.now(), LocalTime.of(7,0)), localDateTime1);
        Assert.assertEquals(LocalDateTime.of(LocalDate.now(), LocalTime.of(21,50)), localDateTime2);
        Assert.assertEquals(LocalDateTime.of(LocalDate.now(), LocalTime.of(21,50)), localDateTime3);
        Assert.assertEquals(LocalDateTime.of(LocalDate.now().plus(1, ChronoUnit.DAYS),
                LocalTime.of(5,0)), localDateTime4);
    }

    @Test
    public void testWrongResult() {
        Assert.assertNull(ORFDateConverter.generateDateForEntry(true, LocalDate.now(), ""));
        Assert.assertNull(ORFDateConverter.generateDateForEntry(true, LocalDate.now(), "test:test"));
        Assert.assertNull(ORFDateConverter.generateDateForEntry(true, LocalDate.now(), "500:5000"));
        Assert.assertNull(ORFDateConverter.generateDateForEntry(true, null, ""));
    }
}
