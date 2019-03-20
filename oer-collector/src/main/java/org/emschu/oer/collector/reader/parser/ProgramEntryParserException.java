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

public class ProgramEntryParserException extends ParserException {
    public ProgramEntryParserException() {
    }

    public ProgramEntryParserException(String message) {
        super(message);
    }

    public ProgramEntryParserException(String message, Throwable cause) {
        super(message, cause);
    }

    public ProgramEntryParserException(Throwable cause) {
        super(cause);
    }

    public ProgramEntryParserException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }
}
