package org.emschu.oer.collector.reader.parser;

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

import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.core.model.Channel;

import javax.transaction.Transactional;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;

@Transactional
public abstract class CustomParser {

    private List<LocalDate> dateRangeList = new ArrayList<>();
    private boolean isTvShowCollectingEnabled;
    private boolean isProgramEntryCollectingEnabled;
    private Channel.AdapterFamily adapterFamily;

    public abstract void run() throws ParserException;
    public abstract void cleanup();

    public List<LocalDate> getDateRangeList() {
        return dateRangeList;
    }

    public void setDateRangeList(List<LocalDate> dateRangeList) {
        this.dateRangeList = dateRangeList;
    }

    public boolean isTvShowCollectingEnabled() {
        return isTvShowCollectingEnabled;
    }

    public void setTvShowCollectingEnabled(boolean tvShowCollectingEnabled) {
        isTvShowCollectingEnabled = tvShowCollectingEnabled;
    }

    public boolean isProgramEntryCollectingEnabled() {
        return isProgramEntryCollectingEnabled;
    }

    public void setProgramEntryCollectingEnabled(boolean programEntryCollectingEnabled) {
        isProgramEntryCollectingEnabled = programEntryCollectingEnabled;
    }

    public Channel.AdapterFamily getAdapterFamily() {
        return adapterFamily;
    }

    public void setAdapterFamily(Channel.AdapterFamily adapterFamily) {
        this.adapterFamily = adapterFamily;
    }
}
