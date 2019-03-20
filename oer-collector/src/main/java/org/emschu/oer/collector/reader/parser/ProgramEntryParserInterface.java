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

import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.collector.reader.ParserException;

import java.time.LocalDate;
import java.util.List;

/**
 * needs to be implemented by any program entry parser.
 */
public interface ProgramEntryParserInterface <T> {
    /**
     * this method is called on every element, you returned in {@link #getElements(Channel, LocalDate)}.
     * use this method to define basic identity fields - e.g. the entries url, technical id ... - of your program entry.
     *
     * @param html
     * @param affectedDay
     * @return
     * @throws ProgramEntryParserException
     */
    public ProgramEntry preProcessItem(T element, LocalDate affectedDay) throws ProgramEntryParserException;

    /**
     * This method is called on each program entry record, right after it was pre-processed.
     *
     * @param programEntry
     * @return
     * @throws ProgramEntryParserException
     */
    public void postProcessItem(ProgramEntry programEntry) throws ProgramEntryParserException;

    /**
     * link stored program entries to other entities here
     *
     * @param programEntry
     */
    public void linkItem(ProgramEntry programEntry);

    /**
     * Filter out the relevant html elements
     *
     * @param channel
     * @param day
     * @return
     * @throws ParserException
     */
    public Iterable<T> getElements(Channel channel, LocalDate day) throws ParserException;

    public void cleanup();

    /**
     * Last action of a single program entry post-process thread
     * use this to free resources
     *
     * @param programEntry
     */
    public void finishEntry(ProgramEntry programEntry);

    void preProcessProgramList(List<ProgramEntry> linkedProgramList);
}
