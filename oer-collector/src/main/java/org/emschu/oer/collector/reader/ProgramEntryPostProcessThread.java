package org.emschu.oer.collector.reader;

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

import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.repository.ProgramEntryRepository;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.logging.Logger;

/**
 * this thread is used for program entries to fetch its data. As this may be time consuming (other api calls...)
 * we use a {@link Runnable} implementation here.
 *
 * <p>
 * suppressing printStacktrace hint in sonar
 */
@SuppressWarnings("squid:S1148")
public class ProgramEntryPostProcessThread implements Runnable {
    private final ProgramEntry programEntry;
    private final ProgramEntryParserInterface programEntryParser;
    private final ProgramEntryRepository programEntryRepository;
    private static final Logger LOG = Logger.getLogger(ProgramEntryPostProcessThread.class.getName());
    private final boolean isDebug;

    /**
     * thread constructor
     *
     * @param programEntryParser
     * @param programEntryRepository
     * @param programEntry
     * @param isDebug bool flag
     */
    public ProgramEntryPostProcessThread(ProgramEntryParserInterface programEntryParser, ProgramEntryRepository programEntryRepository,
                                         ProgramEntry programEntry, boolean isDebug) {
        this.programEntry = programEntry;
        this.programEntryParser = programEntryParser;
        this.programEntryRepository = programEntryRepository;
        this.isDebug = isDebug;
    }

    @Override
    public void run() {
        try {
            programEntryParser.postProcessItem(this.programEntry);
        } catch (Exception e) {
            LOG.warning("Exception occured. Skipping entry: " + e.getMessage());
            LOG.throwing(ProgramEntryPostProcessThread.class.getName(), "run", e);
            return;
        }
        if (this.programEntry.getId() != null) {
            this.programEntry.setUpdatedAt(LocalDateTime.now());
        }
        // calc minutes of program entry finally
        // calc duration in minutes
        if (programEntry.getStartDateTime() != null && programEntry.getEndDateTime() != null) {
            long minutes = ChronoUnit.MINUTES.between(programEntry.getStartDateTime(), programEntry.getEndDateTime());
            programEntry.setDurationInMinutes((int) minutes);
        } else {
            LOG.info("problem with entry: " + programEntry.toString());
            throw new IllegalStateException("missing start or end date in program entry: " + programEntry.getTitle());
        }

        try {
            ProgramEntry entry = programEntryRepository.save(this.programEntry);
            programEntryParser.linkItem(entry);
        } catch (Exception e) {
            LOG.warning("program entry could not be stored: " + e.getMessage());
            if (isDebug) {
                LOG.warning("Program entry:" + this.programEntry);
                LOG.info("Debug mode stacktrace: ");
                e.printStackTrace();
            }
        } finally {
            programEntryParser.finishEntry(this.programEntry);
        }
    }
}
