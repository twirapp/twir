import { MultiSelect } from '@mantine/core';
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';

const GameAliasesCreator = () => {
  const [gameAliases, setGameAliases] = useState<Array<string>>([]);
  const [gameAliasesSearch, setGameAliasesSearch] = useState('');
  const { t } = useTranslation('commands');

  return (
    <>
      <MultiSelect
        label={t('drawer.gameAliases.name')}
        data={gameAliases}
        placeholder={t('drawer.gameAliases.placeholder')!}
        searchable
        creatable
        withinPortal
        getCreateLabel={(query) => `+ Create ${query}`}
        onChange={(data) => {
          setGameAliases(data);
        }}
        searchValue={gameAliasesSearch}
        onSearchChange={setGameAliasesSearch}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ';' || e.key === ',') {
            if (gameAliases.includes(gameAliasesSearch)) return;
            setGameAliases((data) => [...data, gameAliasesSearch]);
            setGameAliasesSearch('');
          }
        }}
      />
    </>
  );
};

export default GameAliasesCreator;
