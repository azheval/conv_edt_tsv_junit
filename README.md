# Converter EDT tsv to junit xml

![build](https://github.com/azheval/conv_edt_tsv_junit/actions/workflows/main.yaml/badge.svg)

Конвертер результатов проверки конфигурации вызовом EDT в junit xml.

## Использование

`./conv_edt_tsv_junit-windows-amd64.exe --settings_file=config.json`

## Настройка

- 'input_file_folder': директория с результатами проверки
- 'output_file_folder': директория с результатами конвертации
- 'skip_errors_file': файл проверки конфигурации, результаты которой нужно пропустить при текущей проверке, значение может быть пустым
- 'skip_categories': категории проверки, которые будут пропущены при конвертации
- 'skip_objects': объекты проверки, которые будут пропущены при конвертации
- 'skip_significance_categories': значимости и категории проверки, которые будут пропущены при конвертации
- 'skip_error_text': ошибки, которые будут пропущены при конвертации

## Конвертация

Строка файла src_file_name.tsv:
```2024-07-17T15:04:48+0300	Ошибка конфигурации		cf		Обработка.ОбменСПорталомСТТ.Форма.ФормаОбработки.Форма.Модуль	строка 1036	Функция 'ПолучитьИмяВременногоФайла' не определена [Web-клиент]```

конвертируется в:
```xml
<testsuite name="src_file_name_Ошибка конфигурации_" timestamp="2025-02-18T15:07:51" time="0" tests="1" errors="0" failures="1" skipped="0">
    <properties></properties>
    <testcase classname="" name="Обработка.ОбменСПорталомСТТ.Форма.ФормаОбработки.Форма.Модуль" time="0.010000">
        <failure message="Ошибка конфигурации; ; " type="">Обработка.ОбменСПорталомСТТ.Форма.ФормаОбработки.Форма.Модуль; строка 1036; Функция &#39;ПолучитьИмяВременногоФайла&#39; не определена [Web-клиент]</failure>
    </testcase>
</testsuite>```

Если такая же строка присутствует в файле, указанном в 'skip_errors_file', или попадет под соответствие одного из фильтров 'skip...', то она будет пропущена.

