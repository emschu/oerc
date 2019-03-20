package org.emschu.oer.collector.reader;

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

import java.util.concurrent.ConcurrentHashMap;

public abstract class AbstractTvShowParser {

    private ConcurrentHashMap<String, String> usedTvShowIds = new ConcurrentHashMap<>();

    protected void registerKey(String tvShowId, String tvShowHash) {
        if (isKeyRegistered(tvShowId)) {
            throw new IllegalStateException(String.format("key '%s' is already registered", tvShowId));
        }
        usedTvShowIds.put(tvShowId, tvShowHash);
    }

    protected boolean isKeyRegistered(String key) {
        return usedTvShowIds.containsKey(key);
    }

    public void clear() {
        this.usedTvShowIds.clear();
    }
}
