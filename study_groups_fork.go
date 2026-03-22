package main

import (
	"log"
)

func (h *APIHandler) markStudyGroupInstallsForkedByDeckIDs(deckIDs ...int64) {
	seen := make(map[int64]struct{}, len(deckIDs))
	for _, deckID := range deckIDs {
		if deckID == 0 {
			continue
		}
		if _, ok := seen[deckID]; ok {
			continue
		}
		seen[deckID] = struct{}{}
		if err := h.store.MarkStudyGroupInstallForkedByDeckID(deckID); err != nil {
			log.Printf("failed to mark study group install forked for deck %d: %v", deckID, err)
		}
	}
}

func (h *APIHandler) markStudyGroupInstallsForkedByNoteType(noteTypeName string) {
	if noteTypeName == "" {
		return
	}
	if err := h.store.MarkStudyGroupInstallsForkedByNoteType(h.collectionID, noteTypeName); err != nil {
		log.Printf("failed to mark study group installs forked for note type %s: %v", noteTypeName, err)
	}
}
