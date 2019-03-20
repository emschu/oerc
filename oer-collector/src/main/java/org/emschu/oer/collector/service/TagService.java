package org.emschu.oer.collector.service;

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

import org.emschu.oer.core.model.Tag;
import org.emschu.oer.core.model.repository.TagRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.transaction.Transactional;
import java.time.LocalDateTime;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class TagService {

    @Autowired
    private TagRepository tagRepository;

    private ConcurrentHashMap<String, Long> storedTagIds = new ConcurrentHashMap<>();

    // till we find the entries in db, store them here.
    // always has to be as small as possible
    private ConcurrentHashMap<String, Tag> tmpStoredTags = new ConcurrentHashMap<>();

    /**
     * @param tagName
     * @return
     */
    @Transactional
    public synchronized Tag getOrCreateTag(String tagName) {
        // "normalize"
        tagName = tagName.trim();
        if (storedTagIds.containsKey(tagName)) {
            Long id = storedTagIds.get(tagName);
            Optional<Tag> tagOptional = tagRepository.findById(id);
            // waiting for jpa commit of lines below, can't be long...
            if (!tagOptional.isPresent() && tmpStoredTags.containsKey(tagName)) {
                return tmpStoredTags.get(tagName);
            }
            return tagOptional.orElseThrow(IllegalStateException::new);
        }
        // use db or instantly create a new one
        Optional<Tag> dbTag = tagRepository.findByTagName(tagName);
        if (dbTag.isPresent()) {
            storedTagIds.put(tagName, dbTag.get().getId());
            return dbTag.get();
        }
        // create new tag
        Tag newTag = new Tag();
        newTag.setCreatedAt(LocalDateTime.now());
        newTag.setTagName(tagName);
        // just before saving, we check, if the concurrent hash map has an entry now..
        newTag = tagRepository.save(newTag);
        storedTagIds.put(newTag.getTagName(), newTag.getId());
        tmpStoredTags.put(newTag.getTagName(), newTag);
        return newTag;
    }

    public void clear() {
        tmpStoredTags.clear();
        storedTagIds.clear();
    }
}
