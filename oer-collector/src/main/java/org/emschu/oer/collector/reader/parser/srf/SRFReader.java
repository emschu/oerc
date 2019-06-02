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

import org.emschu.oer.collector.reader.AbstractReader;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.collector.reader.parser.TvShowParserInterface;
import org.emschu.oer.core.model.Channel;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component("srfReader")
public class SRFReader extends AbstractReader {

    @Autowired
    private ProgramEntryParser programEntryParser;

    @Autowired
    private TvShowParser tvShowParser;

    @Override
    public ProgramEntryParserInterface getProgramEntryParser() {
        return programEntryParser;
    }

    @Override
    public TvShowParserInterface getTvShowParser() {
        return tvShowParser;
    }

    @Override
    public Channel.AdapterFamily getAdapterFamily() {
        return Channel.AdapterFamily.SRF;
    }
}
