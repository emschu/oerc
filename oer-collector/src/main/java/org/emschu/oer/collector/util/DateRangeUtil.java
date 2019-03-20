package org.emschu.oer.collector.util;

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

import java.time.LocalDate;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * class to generate lists of date objects
 */
public class DateRangeUtil {

    /**
     *
     * @param futureMax incl
     * @param pastMax incl
     * @return
     */
    public static final List<LocalDate> generateDateRangeToFetch(int futureMax, int pastMax) {
        List<LocalDate> dateList = new ArrayList<>();
        dateList.add(LocalDate.now());
        // future first
        if (futureMax > 0) {
            for (int i = 1; i <= futureMax; i++) {
                dateList.add(LocalDate.now().plus(i, ChronoUnit.DAYS));
            }
        }
        if (pastMax > 0) {
            for (int j = 1; j <= pastMax; j++) {
                dateList.add(LocalDate.now().minus(j, ChronoUnit.DAYS));
            }
        }
        Collections.sort(dateList);
        return dateList;
    }

    public static final List<LocalDate> dateRangeBetween(LocalDate startDate, LocalDate endDate) {
        List<LocalDate> dateList = new ArrayList<>();
        if (startDate == null) {
            return dateList;
        }
        if (endDate == null) {
            endDate = LocalDate.now();
        }
        long dayCount = ChronoUnit.DAYS.between(startDate, endDate);
        for (int i = 0; i < dayCount; i++) {
            dateList.add(startDate.plus(i, ChronoUnit.DAYS));
        }
        return dateList;
    }
}
