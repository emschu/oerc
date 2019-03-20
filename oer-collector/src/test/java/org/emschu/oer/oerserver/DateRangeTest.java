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

import junit.framework.TestCase;
import org.junit.Assert;
import org.junit.Test;
import org.emschu.oer.collector.util.DateRangeUtil;

import java.time.LocalDate;
import java.util.List;

public class DateRangeTest extends TestCase {

    @Test
    public void testGenerationOfDateRange() {
        List<LocalDate> dateList = DateRangeUtil.generateDateRangeToFetch(0,0);
        Assert.assertNotEquals(0, dateList.size());
        Assert.assertEquals(1, dateList.size()); // = today

        dateList = DateRangeUtil.generateDateRangeToFetch(1, 0);
        Assert.assertEquals(2, dateList.size());

        dateList = DateRangeUtil.generateDateRangeToFetch(0,1);
        Assert.assertEquals(2, dateList.size());

        dateList = DateRangeUtil.generateDateRangeToFetch(1,1);
        Assert.assertEquals(3, dateList.size());
    }
}
