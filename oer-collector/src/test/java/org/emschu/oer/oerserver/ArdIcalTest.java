package org.emschu.oer.oerserver;

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

import org.apache.commons.io.IOUtils;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.emschu.oer.collector.reader.parser.ard.ArdICalParser;
import org.springframework.data.util.Pair;

import java.io.IOException;
import java.time.LocalDateTime;

public class ArdIcalTest {

    private ArdICalParser parser;

    @Before
    public void setup() {
        parser = new ArdICalParser();
    }

    @Test
    public void testExpectedStandardContent() throws IOException {
        String icalContent = IOUtils.toString(
                this.getClass().getResourceAsStream("/other_ical_event.txt"),
                "UTF-8"
        );
        Pair<LocalDateTime, LocalDateTime> dateValues = parser.getRelevantDates(icalContent);
        Assert.assertNotNull(dateValues);
        Assert.assertNotNull(dateValues.getFirst());
        Assert.assertNotNull(dateValues.getSecond());
        Assert.assertEquals(LocalDateTime.of(2019, 2, 25, 9, 35), dateValues.getFirst());
        Assert.assertEquals(LocalDateTime.of(2019, 2, 25, 10, 05), dateValues.getSecond());
    }

    @Test
    public void testEvent1() throws IOException {
        String icalContent = IOUtils.toString(
                this.getClass().getResourceAsStream("/ical_event.txt"),
                "UTF-8"
        );
        Pair<LocalDateTime, LocalDateTime> dateValues = parser.getRelevantDates(icalContent);
        Assert.assertNotNull(dateValues);
        Assert.assertNotNull(dateValues.getFirst());
        Assert.assertNotNull(dateValues.getSecond());
        Assert.assertEquals(LocalDateTime.of(2019, 2, 25, 10, 33), dateValues.getFirst());
        Assert.assertEquals(LocalDateTime.of(2019, 2, 25, 10, 35), dateValues.getSecond());
    }
}
