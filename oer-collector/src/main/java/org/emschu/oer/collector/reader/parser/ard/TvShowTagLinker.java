package org.emschu.oer.collector.reader.parser.ard;

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

import org.hibernate.Hibernate;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.parser.CustomParser;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.Tag;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.core.model.repository.TvShowRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.Set;
import java.util.logging.Logger;

@Component(value = "ardTvShowTagLinker")
public class TvShowTagLinker extends CustomParser {

    private static final Logger LOG = Logger.getLogger(ProgramEntryTagLinker.class.getName());

    @Autowired
    private TvShowRepository tvShowRepository;

    @Override
    public void run() throws ParserException {
        Iterable<TvShow> tvShowStream = tvShowRepository.findAllByAdapterFamily(getAdapterFamily());
        tvShowStream.forEach(tvShow -> {
            Set<Tag> tagSet = new HashSet<>();

            if (!Hibernate.isInitialized(tvShow) || !Hibernate.isInitialized(tvShow.getRelatedProgramEntries())) {
                Hibernate.initialize(tvShow.getRelatedProgramEntries());
            }

            for (ProgramEntry singleRelatedProgramEntry : tvShow.getRelatedProgramEntries()) {
                for (Tag singleTag : singleRelatedProgramEntry.getTags()) {
                    if (!tvShow.getTags().contains(singleTag)) {
                        tagSet.add(singleTag);
                    }
                }
            }
            tvShow.setTags(new ArrayList<>(tagSet));
            tvShowRepository.save(tvShow);

            if (!tagSet.isEmpty()) {
                LOG.info("Detected " + tagSet.size() + " different tags for tv show " + tvShow.getTitle());
            }
        });
    }

    @Override
    public void cleanup() {

    }
}
