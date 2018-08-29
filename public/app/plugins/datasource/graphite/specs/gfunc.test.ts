import gfunc from '../gfunc';

describe('when creating func instance from func names', function() {
  it('should return func instance', function() {
    const func = gfunc.createFuncInstance('sumSeries');
    expect(func).toBeTruthy();
    expect(func.def.name).toEqual('sumSeries');
    expect(func.def.params.length).toEqual(1);
    expect(func.def.params[0].multiple).toEqual(true);
    expect(func.def.defaultParams.length).toEqual(1);
  });

  it('should return func instance with shortName', function() {
    const func = gfunc.createFuncInstance('sum');
    expect(func).toBeTruthy();
  });

  it('should return func instance from funcDef', function() {
    const func = gfunc.createFuncInstance('sum');
    const func2 = gfunc.createFuncInstance(func.def);
    expect(func2).toBeTruthy();
  });

  it('func instance should have text representation', function() {
    const func = gfunc.createFuncInstance('groupByNode');
    func.params[0] = 5;
    func.params[1] = 'avg';
    func.updateText();
    expect(func.text).toEqual('groupByNode(5, avg)');
  });
});

describe('when rendering func instance', function() {
  it('should handle single metric param', function() {
    const func = gfunc.createFuncInstance('sumSeries');
    expect(func.render('hello.metric')).toEqual('sumSeries(hello.metric)');
  });

  it('should include default params if options enable it', function() {
    const func = gfunc.createFuncInstance('scaleToSeconds', {
      withDefaultParams: true,
    });
    expect(func.render('hello')).toEqual('scaleToSeconds(hello, 1)');
  });

  it('should handle int or interval params with number', function() {
    const func = gfunc.createFuncInstance('movingMedian');
    func.params[0] = '5';
    expect(func.render('hello')).toEqual('movingMedian(hello, 5)');
  });

  it('should handle int or interval params with interval string', function() {
    const func = gfunc.createFuncInstance('movingMedian');
    func.params[0] = '5min';
    expect(func.render('hello')).toEqual("movingMedian(hello, '5min')");
  });

  it('should never quote boolean paramater', function() {
    const func = gfunc.createFuncInstance('sortByName');
    func.params[0] = '$natural';
    expect(func.render('hello')).toEqual('sortByName(hello, $natural)');
  });

  it('should never quote int paramater', function() {
    const func = gfunc.createFuncInstance('maximumAbove');
    func.params[0] = '$value';
    expect(func.render('hello')).toEqual('maximumAbove(hello, $value)');
  });

  it('should never quote node paramater', function() {
    const func = gfunc.createFuncInstance('aliasByNode');
    func.params[0] = '$node';
    expect(func.render('hello')).toEqual('aliasByNode(hello, $node)');
  });

  it('should handle metric param and int param and string param', function() {
    const func = gfunc.createFuncInstance('groupByNode');
    func.params[0] = 5;
    func.params[1] = 'avg';
    expect(func.render('hello.metric')).toEqual("groupByNode(hello.metric, 5, 'avg')");
  });

  it('should handle function with no metric param', function() {
    const func = gfunc.createFuncInstance('randomWalk');
    func.params[0] = 'test';
    expect(func.render(undefined)).toEqual("randomWalk('test')");
  });

  it('should handle function multiple series params', function() {
    const func = gfunc.createFuncInstance('asPercent');
    func.params[0] = '#B';
    expect(func.render('#A')).toEqual('asPercent(#A, #B)');
  });
});

describe('when requesting function definitions', function() {
  it('should return function definitions', function() {
    const funcIndex = gfunc.getFuncDefs('1.0');
    expect(Object.keys(funcIndex).length).toBeGreaterThan(8);
  });
});

describe('when updating func param', function() {
  it('should update param value and update text representation', function() {
    const func = gfunc.createFuncInstance('summarize', {
      withDefaultParams: true,
    });
    func.updateParam('1h', 0);
    expect(func.params[0]).toBe('1h');
    expect(func.text).toBe('summarize(1h, sum, false)');
  });

  it('should parse numbers as float', function() {
    const func = gfunc.createFuncInstance('scale');
    func.updateParam('0.001', 0);
    expect(func.params[0]).toBe('0.001');
  });
});

describe('when updating func param with optional second parameter', function() {
  it('should update value and text', function() {
    const func = gfunc.createFuncInstance('aliasByNode');
    func.updateParam('1', 0);
    expect(func.params[0]).toBe('1');
  });

  it('should slit text and put value in second param', function() {
    const func = gfunc.createFuncInstance('aliasByNode');
    func.updateParam('4,-5', 0);
    expect(func.params[0]).toBe('4');
    expect(func.params[1]).toBe('-5');
    expect(func.text).toBe('aliasByNode(4, -5)');
  });

  it('should remove second param when empty string is set', function() {
    const func = gfunc.createFuncInstance('aliasByNode');
    func.updateParam('4,-5', 0);
    func.updateParam('', 1);
    expect(func.params[0]).toBe('4');
    expect(func.params[1]).toBe(undefined);
    expect(func.text).toBe('aliasByNode(4)');
  });
});
